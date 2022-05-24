package clients

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/lobby"
	"Waffle/bancho/misc"
	"Waffle/bancho/packets"
	"Waffle/database"
	"Waffle/helpers"
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// HandleIncoming handles things coming from the osu! client
func (client *Client) HandleIncoming() {
	//make a 4kb Buffer to read stuff
	readBuffer := make([]byte, 4096)

	for client.continueRunning {
		read, readErr := client.connection.Read(readBuffer)

		if readErr != nil {
			//We don't clean up as we may not need to
			return
		}

		go func() {
			misc.StatsRecvLock.Lock()
			misc.StatsBytesRecieved += uint64(read)
			misc.StatsRecvLock.Unlock()
		}()

		//Update last receive time, this is used to check for timeouts
		client.lastReceive = time.Now()

		//Get the bytes that were actually read
		packetBuffer := bytes.NewBuffer(readBuffer[:read])
		//Index into the buffer, so we read every packet that we have
		readIndex := 0

		for readIndex < read {
			read, packet, failedRead := packets.ReadBanchoPacketHeader(packetBuffer)

			readIndex += read

			if failedRead {
				continue
			}

			//Unused packet
			if packet.PacketId == 79 {
				continue
			}

			//Packet data reader, only contains the packet data
			packetDataReader := bytes.NewBuffer(packet.PacketData)

			switch packet.PacketId {
			//The client is informing us about its new status
			case packets.OsuSendUserStatus:
				statusUpdate := packets.ReadStatusUpdate(packetDataReader)

				client.Status = statusUpdate

				client_manager.BroadcastPacket(func(packetQueue chan packets.BanchoPacket) {
					packets.BanchoSendOsuUpdate(packetQueue, client.GetRelevantUserStats(), client.Status)
				})
			//The client is informing us, that it wants to know its own updated stats
			case packets.OsuRequestStatusUpdate:
				//Retrieve stats
				statGetResultOsu, osuStats := database.UserStatsFromDatabase(client.UserData.UserID, 0)
				statGetResultTaiko, taikoStats := database.UserStatsFromDatabase(client.UserData.UserID, 1)
				statGetResultCatch, catchStats := database.UserStatsFromDatabase(client.UserData.UserID, 2)
				statGetResultMania, maniaStats := database.UserStatsFromDatabase(client.UserData.UserID, 3)

				if statGetResultOsu == -1 || statGetResultTaiko == -1 || statGetResultCatch == -1 || statGetResultMania == -1 {
					packets.BanchoSendAnnounce(client.PacketQueue, "A weird server-side fuckup occured, your stats don't exist yet your user does...")
					return
				} else if statGetResultOsu == -2 || statGetResultTaiko == -2 || statGetResultCatch == -2 || statGetResultMania == -2 {
					packets.BanchoSendAnnounce(client.PacketQueue, "A weird server-side fuckup occured, stats could not be loaded...")
					return
				}

				client.OsuStats = osuStats
				client.TaikoStats = taikoStats
				client.CatchStats = catchStats
				client.ManiaStats = maniaStats

				packets.BanchoSendUserPresence(client.PacketQueue, client.UserData, client.GetRelevantUserStats(), client.GetClientTimezone())
				packets.BanchoSendOsuUpdate(client.PacketQueue, client.GetRelevantUserStats(), client.Status)
			//The Client is requesting more information about certain clients
			case packets.OsuUserStatsRequest:
				var listLength int16

				binary.Read(packetDataReader, binary.LittleEndian, &listLength)

				//Read every user ID the client requests
				for i := 0; int16(i) != listLength; i++ {
					var currentId int32
					binary.Read(packetDataReader, binary.LittleEndian, currentId)

					user := client_manager.GetClientById(currentId)

					//If we didn't find the user, simply skip
					if user == nil {
						continue
					}

					//Send information about the client requested
					packets.BanchoSendOsuUpdate(client.PacketQueue, user.GetRelevantUserStats(), user.GetUserStatus())
					break
				}
			//The client is sending a message into a channel
			case packets.OsuSendIrcMessage:
				message := packets.ReadMessage(packetDataReader)

				//Reroute if it's for #multiplayer
				if message.Target == "#multiplayer" {
					if client.currentMultiLobby != nil {
						client.currentMultiLobby.MultiChannel.SendMessage(client, message.Message, message.Target)

						if message.Message[0] == '!' {
							go client.WaffleBotHandleCommand(client, message)
						}
					}
					break
				}

				//Find channel
				channel, exists := client.joinedChannels[message.Target]

				if exists {
					channel.SendMessage(client, message.Message, message.Target)
					database.ChatInsertNewMessage(uint64(client.GetUserId()), message.Target, message.Message)
				}
				//The client is sending a private message to someone
			case packets.OsuSendIrcMessagePrivate:
				message := packets.ReadMessage(packetDataReader)
				//Assign a sender, as the client doesn't seem to send itself as the sender
				message.Sender = client.UserData.Username

				//Find the target
				targetClient := client_manager.GetClientByName(message.Target)

				//If we found the client, only then send a message
				if targetClient != nil {
					packets.BanchoSendIrcMessage(targetClient.GetPacketQueue(), message)

					awayMessage := targetClient.GetAwayMessage()

					//If the user is marked as away, tell the sender
					if awayMessage != "" {
						packets.BanchoSendIrcMessage(client.PacketQueue, packets.Message{
							Sender:  targetClient.GetUserData().Username,
							Message: fmt.Sprintf("/me is away! %s", awayMessage),
							Target:  client.GetUserData().Username,
						})
					}

					database.ChatInsertNewMessage(uint64(client.GetUserId()), strconv.FormatInt(int64(targetClient.GetUserId()), 10), message.Message)
				}
			//The client nicely informs the server that its leaving
			case packets.OsuExit:
				client.CleanupClient("Player Exited")
				return
			//The client informs that it wants to start spectating someone
			case packets.OsuStartSpectating:
				var spectatorId int32

				binary.Read(packetDataReader, binary.LittleEndian, &spectatorId)

				//Find client to spectate
				toSpectate := client_manager.GetClientById(spectatorId)

				//Leave if none is found
				if toSpectate == nil {
					break
				}

				toSpectate.InformSpectatorJoin(client)

				//Stop spectating old client if there is one
				if client.spectatingClient != nil {
					client.spectatingClient.InformSpectatorLeft(client)
				}

				client.spectatingClient = toSpectate
			//The client informs the server that it wants to stop spectating the current user
			case packets.OsuStopSpectating:
				if client.spectatingClient == nil {
					break
				}

				client.spectatingClient.InformSpectatorLeft(client)
				client.spectatingClient = nil
			//The client is sending replay frames for spectators, this is only done if it knows it has spectators
			case packets.OsuSpectateFrames:
				frameBundle := packets.ReadSpectatorFrameBundle(packetDataReader)

				//Send the frames to all spectators
				client.BroadcastToSpectators(func(packetQueue chan packets.BanchoPacket) {
					packets.BanchoSendSpectateFrames(packetQueue, frameBundle)
				})
			//The client informs the server that it's missing the map which the client its spectating is playing
			case packets.OsuCantSpectate:
				if client.spectatingClient != nil {
					client.spectatingClient.InformSpectatorCantSpectate(client)
				}
			//The client informs the server about an error which had occurred
			case packets.OsuErrorReport:
				errorString := string(packets.ReadBanchoString(packetDataReader))

				helpers.Logger.Printf("[Bancho@Handling] %s Encountered an error!! Error Details:\n%s", client.UserData.Username, errorString)
			//This is the response to a BanchoPing
			case packets.OsuPong:
				client.lastReceive = time.Now()
			//The client has joined the lobby
			case packets.OsuLobbyJoin:
				lobby.JoinLobby(client)
				client.isInLobby = true
			//The client has left the lobby
			case packets.OsuLobbyPart:
				lobby.PartLobby(client)
				client.isInLobby = false
			//The client is requesting to join a chat channel
			case packets.OsuChannelJoin:
				channelName := string(packets.ReadBanchoString(packetDataReader))

				channel, exists := chat.GetChannelByName(channelName)

				//If the channel exists, attempt to join
				if exists {
					if channel.Join(client) {
						packets.BanchoSendChannelJoinSuccess(client.PacketQueue, channelName)
						client.joinedChannels[channel.Name] = channel
					} else {
						packets.BanchoSendChannelRevoked(client.PacketQueue, channelName)
					}
				} else {
					packets.BanchoSendChannelRevoked(client.PacketQueue, channelName)
				}
			//The client is requesting to leave a chat channel
			case packets.OsuChannelLeave:
				channelName := string(packets.ReadBanchoString(packetDataReader))

				//Search for the channel
				channel, exists := client.joinedChannels[channelName]

				if exists {
					channel.Leave(client)
					delete(client.joinedChannels, channelName)
				}
			//The client is creating a multiplayer match
			case packets.OsuMatchCreate:
				match := packets.ReadMultiplayerMatch(packetDataReader)

				lobby.CreateNewMultiMatch(match, client)
			//The client is looking to join a multiplayer match
			case packets.OsuMatchJoin:
				matchJoin := packets.ReadMatchJoin(packetDataReader)

				foundMatch := lobby.GetMultiMatchById(uint16(matchJoin.MatchId))

				//Only try joining if one is found
				if foundMatch != nil {
					client.JoinMatch(foundMatch, matchJoin.Password)
				} else {
					packets.BanchoSendMatchJoinFail(client.PacketQueue)
				}
			//The client wants to leave the current multiplayer match
			case packets.OsuMatchPart:
				client.LeaveCurrentMatch()
			//The client wants to change in which multiplayer slot its in
			case packets.OsuMatchChangeSlot:
				if client.currentMultiLobby != nil {
					var newSlot int32

					binary.Read(packetDataReader, binary.LittleEndian, &newSlot)

					client.currentMultiLobby.TryChangeSlot(client, int(newSlot))
				}
			//The client wants to change sides
			case packets.OsuMatchChangeTeam:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.ChangeTeam(client)
				}
			//The client wants to transfer its host status onto someone else
			case packets.OsuMatchTransferHost:
				if client.currentMultiLobby != nil {
					var newHost int32

					binary.Read(packetDataReader, binary.LittleEndian, &newHost)

					client.currentMultiLobby.TransferHost(client, int(newHost))
				}
			//The client informs the server it has pressed the ready button
			case packets.OsuMatchReady:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.ReadyUp(client)
				}
			//The client informs the server it has pressed the not ready button
			case packets.OsuMatchNotReady:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.Unready(client)
				}
			//The client informs the server it has made some changes to the settings of the match
			case packets.OsuMatchChangeSettings:
				if client.currentMultiLobby != nil {
					newMatch := packets.ReadMultiplayerMatch(packetDataReader)
					client.currentMultiLobby.ChangeSettings(client, newMatch)
				}
			//The client informs the server that it has changed the mods in the match
			case packets.OsuMatchChangeMods:
				if client.currentMultiLobby != nil {
					var newMods int32

					binary.Read(packetDataReader, binary.LittleEndian, &newMods)

					client.currentMultiLobby.ChangeMods(client, newMods)
				}
			//The client informs the server that it has tried to lock a slot in the multi lobby
			case packets.OsuMatchLock:
				if client.currentMultiLobby != nil {
					var slot int32

					binary.Read(packetDataReader, binary.LittleEndian, &slot)

					client.currentMultiLobby.LockSlot(client, int(slot))
				}
			//The client informs the server that it's missing the beatmap to be played
			case packets.OsuMatchNoBeatmap:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.InformNoBeatmap(client)
				}
			//The client informs the server that it now has the beatmap that is to be played
			case packets.OsuMatchHasBeatmap:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.InformGotBeatmap(client)
				}
			//The client informs the server that it has completed playing the map
			case packets.OsuMatchComplete:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.InformCompletion(client)
				}
			//The client informs the server that it has loaded into the game successfully
			case packets.OsuMatchLoadComplete:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.InformLoadComplete(client)
				}
			//The client informs the server of its new score
			case packets.OsuMatchScoreUpdate:
				if client.currentMultiLobby != nil {
					scoreFrame := packets.ReadScoreFrame(packetDataReader)

					client.currentMultiLobby.InformScoreUpdate(client, scoreFrame)
				}
			//The client has requested to skip the beginning break
			case packets.OsuMatchSkipRequest:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.InformPressedSkip(client)
				}
			//The client has failed the map
			case packets.OsuMatchFailed:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.InformFailed(client)
				}
			//The client has pressed start game
			case packets.OsuMatchStart:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.StartGame(client)
				}
			//The client is looking to add a friend to their friends list
			case packets.OsuFriendsAdd:
				var friendId int32

				binary.Read(packetDataReader, binary.LittleEndian, &friendId)

				//Append friends list
				client.FriendsList = append(client.FriendsList, database.FriendEntry{
					User1: client.UserData.UserID,
					User2: uint64(friendId),
				})

				//Save in database
				go database.FriendsAddFriend(client.UserData.UserID, uint64(friendId))
			//The client is looking to remove a friend from their friends list
			case packets.OsuFriendsRemove:
				var friendId int32

				binary.Read(packetDataReader, binary.LittleEndian, &friendId)

				for index, value := range client.FriendsList {
					if value.User2 == uint64(friendId) {
						client.FriendsList = append(client.FriendsList[0:index], client.FriendsList[index+1:]...)
					}
				}

				go database.FriendsRemoveFriend(client.UserData.UserID, uint64(friendId))
			//The client is setting their away message
			case packets.OsuSetIrcAwayMessage:
				awayMessage := packets.ReadMessage(packetDataReader)

				client.awayMessage = awayMessage.Message

				//Setting it empty resets it
				if awayMessage.Message == "" {
					packets.BanchoSendIrcMessage(client.PacketQueue, packets.Message{
						Sender:  "WaffleBot",
						Message: "You're no longer marked as away!",
						Target:  client.UserData.Username,
					})
				} else {
					packets.BanchoSendIrcMessage(client.PacketQueue, packets.Message{
						Sender:  "WaffleBot",
						Message: fmt.Sprintf("You're now marked as away: %s", awayMessage.Message),
						Target:  client.UserData.Username,
					})
				}
			case packets.OsuBeatmapInfoRequest:
				infoRequest := packets.ReadBeatmapInfoRequest(packetDataReader)

				go func() {
					infoReply := packets.BeatmapInfoReply{}

					//Initially store the user ids for the prepared statement
					queryArguments := []interface{}{
						client.OsuStats.UserID, client.OsuStats.UserID, client.OsuStats.UserID,
					}

					//will store the prepared statement question marks for the filenames
					questionMarks := ""

					//edge case for no filenames, it will still have at least 1 filename, even though itll be empty
					if len(infoRequest.Filenames) == 0 {
						questionMarks = "?"
						queryArguments = append(queryArguments, "")
					} else {
						//for every filename add a question mark
						for i := 0; i != len(infoRequest.Filenames); i++ {
							questionMarks += "?, "
						}
					}

					//trim off the last comma to not cause massive issues
					questionMarks = strings.TrimSuffix(questionMarks, ", ")

					//will store the beatmap ids for the sql
					beatmapIds := ""

					//edge case for no beatmap ids, it will still have at least 1 beatmap id, even though itll be 0
					if len(infoRequest.BeatmapIds) == 0 {
						beatmapIds = "0"
					} else {
						//add every beatmap id
						for i := 0; i != len(infoRequest.BeatmapIds); i++ {
							beatmapIds += string(infoRequest.BeatmapIds[i]) + ", "
						}
					}

					//trim off comma to not have a extra one
					beatmapIds = strings.TrimSuffix(beatmapIds, ", ")

					sqlString := `
SELECT				
	result.beatmap_id,				
	result.beatmapset_id,
	result.filename,
	result.beatmap_md5,
	result.ranking_status,
	result.final_osu_ranking AS 'osu_ranking',
	result.final_taiko_ranking AS 'taiko_ranking',
	result.final_catch_ranking AS 'catch_ranking'
FROM (
	SELECT beatmaps.beatmap_id, 
		   beatmaps.beatmapset_id, 
		   beatmaps.filename, 
		   beatmaps.beatmap_md5, 
		   beatmaps.ranking_status, 
		   osuResult.ranking AS 'osu_ranking', 
		   osuResult.user_id AS 'osu_user_id', 
		   taikoResult.ranking AS 'taiko_ranking',
		   taikoResult.user_id AS 'taiko_user_id', 
		   catchResult.ranking AS 'catch_ranking', 
		   catchResult.user_id AS 'catch_user_id',
		CASE WHEN osuResult.ranking IS NULL THEN 'N' ELSE osuResult.ranking END AS 'final_osu_ranking',
		CASE WHEN taikoResult.ranking IS NULL THEN 'N' ELSE taikoResult.ranking END AS 'final_taiko_ranking', 
		CASE WHEN catchResult.ranking IS NULL THEN 'N' ELSE catchResult.ranking END AS 'final_catch_ranking'
			FROM waffle.beatmaps 
		LEFT JOIN scores osuResult ON osuResult.beatmap_id = beatmaps.beatmap_id AND osuResult.mapset_best = 1 AND osuResult.playmode = 0 AND osuResult.user_id = ? 
		LEFT JOIN scores taikoResult ON taikoResult.beatmap_id = beatmaps.beatmap_id AND taikoResult.mapset_best = 1 AND taikoResult.playmode = 1 AND taikoResult.user_id = ? 
		LEFT JOIN scores catchResult ON catchResult.beatmap_id = beatmaps.beatmap_id AND catchResult.mapset_best = 1 AND catchResult.playmode = 2 AND catchResult.user_id = ? 
	WHERE beatmaps.filename IN ( %s ) 
	OR beatmaps.beatmap_id IN ( %s )
) result
`
					//the absolutely gigantic sql
					sql := fmt.Sprintf(sqlString, questionMarks, beatmapIds)

					//add the filenames as query arguments
					for i := 0; i != len(infoRequest.Filenames); i++ {
						queryArguments = append(queryArguments, infoRequest.Filenames[i])
					}

					//query
					databaseQuery, databaseQueryErr := database.Database.Query(sql, queryArguments...)

					//momentarily nonexistant error handling
					if databaseQueryErr != nil {
						if databaseQuery != nil {
							databaseQuery.Close()
						}
					}

					if databaseQuery != nil {
						//for every database result
						for databaseQuery.Next() {
							beatmapInfo := packets.BeatmapInfo{}

							beatmapInfo.InfoId = -1

							var osuRank, taikoRank, catchRank string
							var osuFilename string

							databaseQuery.Scan(&beatmapInfo.BeatmapId, &beatmapInfo.BeatmapSetId, &osuFilename, &beatmapInfo.BeatmapChecksum, &beatmapInfo.Ranked, &osuRank, &taikoRank, &catchRank)

							//convert string rank to peppys enum ranking
							switch osuRank {
							case "XH":
								beatmapInfo.OsuRank = 0
							case "SH":
								beatmapInfo.OsuRank = 1
							case "X":
								beatmapInfo.OsuRank = 2
							case "S":
								beatmapInfo.OsuRank = 3
							case "A":
								beatmapInfo.OsuRank = 4
							case "B":
								beatmapInfo.OsuRank = 5
							case "C":
								beatmapInfo.OsuRank = 6
							case "D":
								beatmapInfo.OsuRank = 7
							case "F":
								beatmapInfo.OsuRank = 8
							case "N":
								beatmapInfo.OsuRank = 9
							}

							//convert string rank to peppys enum ranking
							switch taikoRank {
							case "XH":
								beatmapInfo.TaikoRank = 0
							case "SH":
								beatmapInfo.TaikoRank = 1
							case "X":
								beatmapInfo.TaikoRank = 2
							case "S":
								beatmapInfo.TaikoRank = 3
							case "A":
								beatmapInfo.TaikoRank = 4
							case "B":
								beatmapInfo.TaikoRank = 5
							case "C":
								beatmapInfo.TaikoRank = 6
							case "D":
								beatmapInfo.TaikoRank = 7
							case "F":
								beatmapInfo.TaikoRank = 8
							case "N":
								beatmapInfo.TaikoRank = 9
							}

							//convert string rank to peppys enum ranking
							switch catchRank {
							case "XH":
								beatmapInfo.CatchRank = 0
							case "SH":
								beatmapInfo.CatchRank = 1
							case "X":
								beatmapInfo.CatchRank = 2
							case "S":
								beatmapInfo.CatchRank = 3
							case "A":
								beatmapInfo.CatchRank = 4
							case "B":
								beatmapInfo.CatchRank = 5
							case "C":
								beatmapInfo.CatchRank = 6
							case "D":
								beatmapInfo.CatchRank = 7
							case "F":
								beatmapInfo.CatchRank = 8
							case "N":
								beatmapInfo.CatchRank = 9
							}

							//will store the index of the info request the client gave us
							infoPosition := int16(-1)

							//find it
							for k := 0; k != len(infoRequest.Filenames); k++ {
								if infoRequest.Filenames[k] == osuFilename {
									infoPosition = int16(k)
								}
							}

							beatmapInfo.InfoId = infoPosition

							//append to the reply list
							infoReply.BeatmapInfos = append(infoReply.BeatmapInfos, beatmapInfo)
						}

						//send off
						packets.BanchoSendBeatmapInfoReply(client.PacketQueue, infoReply)

						//make sure to close the connection
						databaseQuery.Close()
					}
				}()
			default:
				helpers.Logger.Printf("[Bancho@Handling] %s: Got %s, of Size: %d\n", client.UserData.Username, packets.GetPacketName(packet.PacketId), packet.PacketSize)
			}
		}
	}
}

// SendOutgoing is looping over the packet queue and waiting for new packets, and sends them off as they come in
func (client *Client) SendOutgoing() {
	for packet := range client.PacketQueue {
		if packet.PacketId != 8 {
			helpers.Logger.Printf("[Bancho@Handling] Sending %s to %s\n", packets.GetPacketName(packet.PacketId), client.UserData.Username)
		}

		sendBytes := packet.GetBytes()

		go func() {
			misc.StatsSendLock.Lock()
			misc.StatsBytesSent += uint64(len(sendBytes))
			misc.StatsSendLock.Unlock()
		}()

		client.connection.Write(packet.GetBytes())
	}
}

// MaintainClient is looping every second, sending out pings and handles timeouts
func (client *Client) MaintainClient() {
	for client.continueRunning {
		lastReceiveUnix := client.lastReceive.Unix()
		lastPingUnix := client.lastPing.Unix()
		unixNow := time.Now().Unix()

		if lastReceiveUnix+ReceiveTimeout <= unixNow {
			client.CleanupClient("Client Timed out.")

			client.continueRunning = false
		}

		if lastPingUnix+PingTimeout <= unixNow {
			packets.BanchoSendPing(client.PacketQueue)

			client.lastPing = time.Now()
		}

		time.Sleep(time.Second)
	}

	//We close in MaintainClient instead of in CleanupClient to avoid possible double closes, causing panics
	helpers.Logger.Printf("[Bancho@Handling] Closed %s's Packet Queue", client.UserData.Username)
	close(client.PacketQueue)
}
