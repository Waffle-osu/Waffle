package b1815

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/clients"
	"Waffle/bancho/lobby"
	"Waffle/bancho/misc"
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/database"
	"Waffle/helpers"
	"Waffle/helpers/serialization"
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"time"
)

// HandleIncoming handles things coming from the osu! client
func (client *Client) HandleIncoming() {
	//make a 4kb Buffer to read stuff
	readBuffer := make([]byte, 32768)

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
			read, packet, failedRead := serialization.ReadBanchoPacketHeader(packetBuffer)

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
			case serialization.OsuSendUserStatus:
				statusUpdate := base_packet_structures.ReadStatusUpdate(packetDataReader)

				client.Status = statusUpdate

				client_manager.BroadcastPacketOsu(func(broadcastClient client_manager.WaffleClient) {
					broadcastClient.BanchoOsuUpdate(client.GetRelevantUserStats(), client.GetUserStatus())
				})

			//The client is informing us, that it wants to know its own updated stats
			case serialization.OsuRequestStatusUpdate:
				//Retrieve stats
				statGetResultOsu, osuStats := database.UserStatsFromDatabase(client.UserData.UserID, 0)
				statGetResultTaiko, taikoStats := database.UserStatsFromDatabase(client.UserData.UserID, 1)
				statGetResultCatch, catchStats := database.UserStatsFromDatabase(client.UserData.UserID, 2)

				if statGetResultOsu == -1 || statGetResultTaiko == -1 || statGetResultCatch == -1 {
					client.BanchoAnnounce("A weird server-side fuckup occured, your stats don't exist yet your user does...")
					return
				} else if statGetResultOsu == -2 || statGetResultTaiko == -2 || statGetResultCatch == -2 {
					client.BanchoAnnounce("A weird server-side fuckup occured, stats could not be loaded...")
					return
				}

				client.OsuStats = osuStats
				client.TaikoStats = taikoStats
				client.CatchStats = catchStats

				client.BanchoPresence(client.UserData, client.GetRelevantUserStats(), client.GetClientTimezone())
				client.BanchoOsuUpdate(client.GetRelevantUserStats(), client.Status)
			//The Client is requesting more information about certain clients
			case serialization.OsuUserStatsRequest:
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
					client.BanchoOsuUpdate(user.GetRelevantUserStats(), client.GetUserStatus())
					break
				}
			//The client is sending a message into a channel
			case serialization.OsuSendIrcMessage:
				if time.Now().Unix() < client.silencedUntil {
					client.SendChatMessage("WaffleBot", fmt.Sprintf("You're silenced for at least %d seconds!", client.silencedUntil-time.Now().Unix()), client.UserData.Username)
				} else {

					message := base_packet_structures.ReadMessage(packetDataReader)

					//Reroute if it's for #multiplayer
					if message.Target == "#multiplayer" {
						if client.currentMultiLobby != nil {
							client.currentMultiLobby.MultiChannel.SendMessage(client, message.Message, message.Target)

							if message.Message[0] == '!' {
								//TODO: wafflebot
								go clients.WaffleBotInstance.WaffleBotHandleCommand(client, message)
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
				}
				//The client is sending a private message to someone
			case serialization.OsuSendIrcMessagePrivate:
				if time.Now().Unix() < client.silencedUntil {
					client.SendChatMessage("WaffleBot", fmt.Sprintf("You're silenced for at least %d seconds!", client.silencedUntil-time.Now().Unix()), client.UserData.Username)
				} else {
					message := base_packet_structures.ReadMessage(packetDataReader)
					//Assign a sender, as the client doesn't seem to send itself as the sender
					message.Sender = client.UserData.Username

					//Find the target
					targetClient := client_manager.GetClientByName(message.Target)

					//If we found the client, only then send a message
					if targetClient != nil {
						targetClient.BanchoIrcMessage(message)

						awayMessage := targetClient.GetAwayMessage()

						//If the user is marked as away, tell the sender
						if awayMessage != "" {
							client.BanchoIrcMessage(base_packet_structures.Message{
								Sender:  targetClient.GetUserData().Username,
								Message: fmt.Sprintf("/me is away! %s", awayMessage),
								Target:  client.GetUserData().Username,
							})
						}

						database.ChatInsertNewMessage(uint64(client.GetUserId()), strconv.FormatInt(int64(targetClient.GetUserId()), 10), message.Message)
					}
				}
			//The client nicely informs the server that its leaving
			case serialization.OsuExit:
				client.CleanupClient("Player Exited")
				return
			//The client informs that it wants to start spectating someone
			case serialization.OsuStartSpectating:
				var spectatorId int32

				binary.Read(packetDataReader, binary.LittleEndian, &spectatorId)

				//Find client to spectate
				toSpectate := client_manager.GetClientById(spectatorId)

				//Leave if none is found
				if toSpectate == nil {
					break
				}

				toSpectate.BanchoSpectatorJoined(client.GetUserId())

				//Stop spectating old client if there is one
				if client.spectatingClient != nil {
					client.spectatingClient.BanchoSpectatorLeft(client.GetUserId())
				}

				client.spectatingClient = toSpectate
			//The client informs the server that it wants to stop spectating the current user
			case serialization.OsuStopSpectating:
				if client.spectatingClient == nil {
					break
				}

				client.spectatingClient.BanchoSpectatorLeft(client.GetUserId())
				client.spectatingClient = nil
			//The client is sending replay frames for spectators, this is only done if it knows it has spectators
			case serialization.OsuSpectateFrames:
				frameBundle := base_packet_structures.ReadSpectatorFrameBundle(packetDataReader)

				//Send the frames to all spectators
				client.BroadcastToSpectators(func(client client_manager.WaffleClient) {
					client.BanchoSpectateFrames(frameBundle)
				})
			//The client informs the server that it's missing the map which the client its spectating is playing
			case serialization.OsuCantSpectate:
				if client.spectatingClient != nil {
					client.spectatingClient.BanchoSpectatorCantSpectate(client.GetUserId())
				}
			//The client informs the server about an error which had occurred
			case serialization.OsuErrorReport:
				errorString := string(serialization.ReadBanchoString(packetDataReader))

				helpers.Logger.Printf("[Bancho@Handling] %s Encountered an error!! Error Details:\n%s", client.UserData.Username, errorString)
			//This is the response to a BanchoPing
			case serialization.OsuPong:
				client.lastReceive = time.Now()
			//The client has joined the lobby
			case serialization.OsuLobbyJoin:
				lobby.JoinLobby(client)
				client.isInLobby = true
			//The client has left the lobby
			case serialization.OsuLobbyPart:
				lobby.PartLobby(client)
				client.isInLobby = false
			//The client is requesting to join a chat channel
			case serialization.OsuChannelJoin:
				channelName := string(serialization.ReadBanchoString(packetDataReader))

				channel, exists := chat.GetChannelByName(channelName)

				//If the channel exists, attempt to join
				if exists {
					if channel.Join(client) {
						client.BanchoChannelJoinSuccess(channelName)

						client.joinedChannels[channel.Name] = channel
					} else {
						client.BanchoChannelRevoked(channelName)
					}
				} else {
					client.BanchoChannelRevoked(channelName)
				}
			//The client is requesting to leave a chat channel
			case serialization.OsuChannelLeave:
				channelName := string(serialization.ReadBanchoString(packetDataReader))

				//Search for the channel
				channel, exists := client.joinedChannels[channelName]

				if exists {
					channel.Leave(client)
					delete(client.joinedChannels, channelName)
				}
			//The client is creating a multiplayer match
			case serialization.OsuMatchCreate:
				match := base_packet_structures.ReadMultiplayerMatch(packetDataReader)

				lobby.CreateNewMultiMatch(match, client)
			//The client is looking to join a multiplayer match
			case serialization.OsuMatchJoin:
				matchJoin := base_packet_structures.ReadMatchJoin(packetDataReader)

				foundMatch := lobby.GetMultiMatchById(uint16(matchJoin.MatchId))

				//Only try joining if one is found
				if foundMatch != nil {
					client.JoinMatch(foundMatch, matchJoin.Password)
				} else {
					client.BanchoMatchJoinFail()
				}
			//The client wants to leave the current multiplayer match
			case serialization.OsuMatchPart:
				client.LeaveCurrentMatch()
			//The client wants to change in which multiplayer slot its in
			case serialization.OsuMatchChangeSlot:
				if client.currentMultiLobby != nil {
					var newSlot int32

					binary.Read(packetDataReader, binary.LittleEndian, &newSlot)

					client.currentMultiLobby.TryChangeSlot(client, int(newSlot))
				}
			//The client wants to change sides
			case serialization.OsuMatchChangeTeam:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.ChangeTeam(client)
				}
			//The client wants to transfer its host status onto someone else
			case serialization.OsuMatchTransferHost:
				if client.currentMultiLobby != nil {
					var newHost int32

					binary.Read(packetDataReader, binary.LittleEndian, &newHost)

					client.currentMultiLobby.TransferHost(client, int(newHost))
				}
			//The client informs the server it has pressed the ready button
			case serialization.OsuMatchReady:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.ReadyUp(client)
				}
			//The client informs the server it has pressed the not ready button
			case serialization.OsuMatchNotReady:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.Unready(client)
				}
			//The client informs the server it has made some changes to the settings of the match
			case serialization.OsuMatchChangeSettings:
				if client.currentMultiLobby != nil {
					newMatch := base_packet_structures.ReadMultiplayerMatch(packetDataReader)
					client.currentMultiLobby.ChangeSettings(client, newMatch)
				}
			//The client informs the server that it has changed the mods in the match
			case serialization.OsuMatchChangeMods:
				if client.currentMultiLobby != nil {
					var newMods int32

					binary.Read(packetDataReader, binary.LittleEndian, &newMods)

					client.currentMultiLobby.ChangeMods(client, newMods)
				}
			//The client informs the server that it has tried to lock a slot in the multi lobby
			case serialization.OsuMatchLock:
				if client.currentMultiLobby != nil {
					var slot int32

					binary.Read(packetDataReader, binary.LittleEndian, &slot)

					client.currentMultiLobby.LockSlot(client, int(slot))
				}
			//The client informs the server that it's missing the beatmap to be played
			case serialization.OsuMatchNoBeatmap:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.InformNoBeatmap(client)
				}
			//The client informs the server that it now has the beatmap that is to be played
			case serialization.OsuMatchHasBeatmap:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.InformGotBeatmap(client)
				}
			//The client informs the server that it has completed playing the map
			case serialization.OsuMatchComplete:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.InformCompletion(client)
				}
			//The client informs the server that it has loaded into the game successfully
			case serialization.OsuMatchLoadComplete:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.InformLoadComplete(client)
				}
			//The client informs the server of its new score
			case serialization.OsuMatchScoreUpdate:
				if client.currentMultiLobby != nil {
					scoreFrame := base_packet_structures.ReadScoreFrame(packetDataReader)

					client.currentMultiLobby.InformScoreUpdate(client, scoreFrame)
				}
			//The client has requested to skip the beginning break
			case serialization.OsuMatchSkipRequest:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.InformPressedSkip(client)
				}
			//The client has failed the map
			case serialization.OsuMatchFailed:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.InformFailed(client)
				}
			//The client has pressed start game
			case serialization.OsuMatchStart:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.StartGame(client)
				}
			//The client is looking to add a friend to their friends list
			case serialization.OsuFriendsAdd:
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
			case serialization.OsuFriendsRemove:
				var friendId int32

				binary.Read(packetDataReader, binary.LittleEndian, &friendId)

				for index, value := range client.FriendsList {
					if value.User2 == uint64(friendId) {
						client.FriendsList = append(client.FriendsList[0:index], client.FriendsList[index+1:]...)
					}
				}

				go database.FriendsRemoveFriend(client.UserData.UserID, uint64(friendId))
			//The client is setting their away message
			case serialization.OsuSetIrcAwayMessage:
				awayMessage := base_packet_structures.ReadMessage(packetDataReader)

				client.awayMessage = awayMessage.Message

				//Setting it empty resets it
				if awayMessage.Message == "" {
					client.BanchoIrcMessage(base_packet_structures.Message{
						Sender:  "WaffleBot",
						Message: "You're no longer marked as away!",
						Target:  client.UserData.Username,
					})
				} else {
					client.BanchoIrcMessage(base_packet_structures.Message{
						Sender:  "WaffleBot",
						Message: fmt.Sprintf("You're now marked as away: %s", awayMessage.Message),
						Target:  client.UserData.Username,
					})
				}
			case serialization.OsuBeatmapInfoRequest:
				infoRequest := base_packet_structures.ReadBeatmapInfoRequest(packetDataReader)

				client.HandleBeatmapInfoRequest(infoRequest)
			default:
				helpers.Logger.Printf("[Bancho@Handling] %s: Got %s, of Size: %d\n", client.UserData.Username, serialization.GetPacketName(packet.PacketId), packet.PacketSize)
			}
		}
	}
}

// SendOutgoing is looping over the packet queue and waiting for new packets, and sends them off as they come in
func (client *Client) SendOutgoing() {
	for packet := range client.PacketQueue {
		sendBytes := len(packet)

		go func() {
			misc.StatsSendLock.Lock()
			misc.StatsBytesSent += uint64(sendBytes)
			misc.StatsSendLock.Unlock()
		}()

		client.connection.Write(packet)
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
			client.BanchoPing()

			client.lastPing = time.Now()
		}

		time.Sleep(time.Second)
	}

	//We close in MaintainClient instead of in CleanupClient to avoid possible double closes, causing panics
	helpers.Logger.Printf("[Bancho@Handling] Closed %s's Packet Queue", client.UserData.Username)
	close(client.PacketQueue)
}
