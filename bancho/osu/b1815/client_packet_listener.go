package b1815

import (
	"Waffle/bancho/bot"
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/lobby"
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/bancho/spectator"
	"Waffle/database"
	"Waffle/helpers"
	"Waffle/helpers/packets"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (client *Client) handlePackets(packetChannel chan packets.BanchoPacket, ctx context.Context) {
	for {
		select {
		case packet := <-packetChannel:
			//Packet data reader, only contains the packet data
			packetDataReader := bytes.NewBuffer(packet.PacketData)

			switch packet.PacketId {
			//The client is informing us about its new status
			case packets.OsuSendUserStatus:
				statusUpdate := packets.Read[base_packet_structures.StatusUpdate](packetDataReader)

				client.Status = statusUpdate

				client_manager.BroadcastPacket(func(broadcastClient client_manager.WaffleClient) {
					broadcastClient.BanchoOsuUpdate(client.GetRelevantUserStats(), client.GetUserStatus())
				})

			//The client is informing us, that it wants to know its own updated stats
			case packets.OsuRequestStatusUpdate:
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
			case packets.OsuUserStatsRequest:
				var listLength int16

				binary.Read(packetDataReader, binary.LittleEndian, &listLength)

				//Read every user ID the client requests
				for i := 0; int16(i) != listLength; i++ {
					var currentId int32
					binary.Read(packetDataReader, binary.LittleEndian, currentId)

					user := client_manager.ClientManager.GetClientById(currentId)

					//If we didn't find the user, simply skip
					if user == nil {
						continue
					}

					//Send information about the client requested
					client.BanchoOsuUpdate(user.GetRelevantUserStats(), client.GetUserStatus())
					break
				}
			//The client is sending a message into a channel
			case packets.OsuSendIrcMessage:
				if time.Now().Unix() < int64(client.UserData.SilencedUntil) {
					client.SendChatMessage("WaffleBot", fmt.Sprintf("You're silenced for at least %d seconds!", int64(client.UserData.SilencedUntil)-time.Now().Unix()), client.UserData.Username)
				} else {
					message := packets.Read[base_packet_structures.Message](packetDataReader)

					//Channel into which the message gets sent,
					//Aswell as where all the WaffleBot/Lobby command
					//Responses get sent to aswell
					var sendChannel *chat.Channel

					//Command outputs
					var returnMessages []string

					//Commands start with !
					if strings.HasPrefix(message.Message, "!") {
						//mp commands take a different route cuz they have to take in LobbyClient
						//instead of just WaffleClient
						if strings.HasPrefix(message.Message, "!mp") {
							returnMessages = lobby.LobbyHandleCommandMultiplayer(client, message.Message)
						} else {
							returnMessages = bot.WaffleBotInstance.WaffleBotHandleCommand(client, message)
						}
					}

					//Reroute for multiplayer
					if message.Target == "#multiplayer" {
						if client.currentMultiLobby != nil {
							sendChannel = client.currentMultiLobby.MultiChannel
						} else {
							client.SendChatMessage("WaffleBot", "Cannot send messages to #multiplayer while not in a Multiplayer Lobby!", "WaffleBot")

							break
						}
					} else {
						//Find channel
						channel, exists := client.joinedChannels[message.Target]

						if !exists {
							client.SendChatMessage("WaffleBot", fmt.Sprintf("No channel found called %s!", message.Target), "WaffleBot")

							break
						}

						sendChannel = channel
					}

					//Send message to appropriate channel
					sendChannel.SendMessage(client, message.Message, message.Target)

					//Follow up with command responses, if any
					for _, content := range returnMessages {
						sendChannel.SendMessage(bot.WaffleBotInstance, content, message.Target)
					}
				}
				//The client is sending a private message to someone
			case packets.OsuSendIrcMessagePrivate:
				if time.Now().Unix() < int64(client.UserData.SilencedUntil) {
					client.SendChatMessage("WaffleBot", fmt.Sprintf("You're silenced for at least %d seconds!", int64(client.UserData.SilencedUntil)-time.Now().Unix()), client.UserData.Username)
				} else {
					message := packets.Read[base_packet_structures.Message](packetDataReader)
					//Assign a sender, as the client doesn't seem to send itself as the sender
					message.Sender = client.UserData.Username

					//Find the target
					targetClient := client_manager.ClientManager.GetClientByName(message.Target)

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
			case packets.OsuExit:
				client.CleanupClient("Player Exited")
				return
			//The client informs that it wants to start spectating someone
			case packets.OsuStartSpectating:
				var spectatorId int32

				binary.Read(packetDataReader, binary.LittleEndian, &spectatorId)

				//Find client to spectate
				toSpectate := spectator.ClientManager.GetClientById(spectatorId)

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
			case packets.OsuStopSpectating:
				if client.spectatingClient == nil {
					break
				}

				client.spectatingClient.BanchoSpectatorLeft(client.GetUserId())
				client.spectatingClient = nil
			//The client is sending replay frames for spectators, this is only done if it knows it has spectators
			case packets.OsuSpectateFrames:
				frameBundle := packets.Read[base_packet_structures.SpectatorFrameBundle](packetDataReader)

				//Send the frames to all spectators
				client.BroadcastToSpectators(func(client spectator.SpectatorClient) {
					client.BanchoSpectateFrames(frameBundle)
				})
			//The client informs the server that it's missing the map which the client its spectating is playing
			case packets.OsuCantSpectate:
				if client.spectatingClient != nil {
					client.spectatingClient.BanchoSpectatorCantSpectate(client.GetUserId())
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
				channelName := packets.Read[string](packetDataReader)

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
				match := packets.Read[base_packet_structures.MultiplayerMatch](packetDataReader)

				lobby.CreateNewMultiMatch(match, client, false)
			//The client is looking to join a multiplayer match
			case packets.OsuMatchJoin:
				matchJoin := packets.Read[base_packet_structures.MatchJoin](packetDataReader)

				foundMatch := lobby.GetMultiMatchById(uint16(matchJoin.MatchId))

				//Only try joining if one is found
				if foundMatch != nil {
					client.JoinMatch(foundMatch, matchJoin.Password)
				} else {
					client.BanchoMatchJoinFail()
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
					newMatch := packets.Read[base_packet_structures.MultiplayerMatch](packetDataReader)
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
					scoreFrame := packets.Read[base_packet_structures.ScoreFrame](packetDataReader)

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
				awayMessage := packets.Read[base_packet_structures.Message](packetDataReader)

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
			case packets.OsuBeatmapInfoRequest:
				infoRequest := packets.Read[base_packet_structures.BeatmapInfoRequest](packetDataReader)

				client.HandleBeatmapInfoRequest(infoRequest)
			default:
				helpers.Logger.Printf("[Bancho@Handling] %s: Got %s, of Size: %d\n", client.UserData.Username, packets.GetPacketName(packet.PacketId), packet.PacketSize)
			}
		case <-ctx.Done():
			return
		}
	}
}
