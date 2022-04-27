package clients

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/database"
	"Waffle/bancho/lobby"
	"Waffle/bancho/packets"
	"bytes"
	"encoding/binary"
	"fmt"
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

		//Update last receive time, this is used to check for timeouts
		client.lastReceive = time.Now()

		//Get the bytes that were actually read
		packetBuffer := bytes.NewBuffer(readBuffer[:read])
		//Index into the buffer, so we read every packet that we have
		readIndex := 0

		for readIndex < read {
			read, packet, failedRead := packets.ReadBanchoPacketHeader(packetBuffer)

			readIndex += read

			if failedRead == true {
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
					packets.BanchoSendOsuUpdate(packetQueue, client.OsuStats, client.Status)
				})
				break
			//The client is informing us, that it wants to know its own stats
			case packets.OsuRequestStatusUpdate:
				packets.BanchoSendUserPresence(client.PacketQueue, client.UserData, client.OsuStats, client.GetClientTimezone())
				packets.BanchoSendOsuUpdate(client.PacketQueue, client.GetRelevantUserStats(), client.Status)
				break
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
					client.currentMultiLobby.MultiChannel.SendMessage(client, message.Message, message.Target)
					break
				}

				//Find channel
				channel, exists := client.joinedChannels[message.Target]

				if exists {
					channel.SendMessage(client, message.Message, message.Target)
				}
				break
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
				}
				break
			//The client nicely informs the server that its leaving
			case packets.OsuExit:
				client.CleanupClient()
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
				break
			//The client informs the server that it wants to stop spectating the current user
			case packets.OsuStopSpectating:
				if client.spectatingClient == nil {
					break
				}

				client.spectatingClient.InformSpectatorLeft(client)
				client.spectatingClient = nil
				break
			//The client is sending replay frames for spectators, this is only done if it knows it has spectators
			case packets.OsuSpectateFrames:
				frameBundle := packets.ReadSpectatorFrameBundle(packetDataReader)

				//Send the frames to all spectators
				client.BroadcastToSpectators(func(packetQueue chan packets.BanchoPacket) {
					packets.BanchoSendSpectateFrames(packetQueue, frameBundle)
				})
				break
			//The client informs the server that it's missing the map which the client its spectating is playing
			case packets.OsuCantSpectate:
				if client.spectatingClient != nil {
					client.spectatingClient.InformSpectatorCantSpectate(client)
				}
				break
			//The client informs the server about an error which had occurred
			case packets.OsuErrorReport:
				errorString := string(packets.ReadBanchoString(packetDataReader))

				fmt.Printf("%s Encountered an error!! Error Details:\n%s", client.UserData.Username, errorString)
				break
			//This is the response to a BanchoPing
			case packets.OsuPong:
				client.lastReceive = time.Now()
				fmt.Printf("%s: Got Ping from %s\n", time.Now().Format(time.ANSIC), client.UserData.Username)
				break
			//The client has joined the lobby
			case packets.OsuLobbyJoin:
				lobby.JoinLobby(client)
				client.isInLobby = true
				break
			//The client has left the lobby
			case packets.OsuLobbyPart:
				lobby.PartLobby(client)
				client.isInLobby = false
				break
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
				break
			//The client is requesting to leave a chat channel
			case packets.OsuChannelLeave:
				channelName := string(packets.ReadBanchoString(packetDataReader))

				//Search for the channel
				channel, exists := client.joinedChannels[channelName]

				if exists {
					channel.Leave(client)
					delete(client.joinedChannels, channelName)
				}
				break
			//The client is creating a multiplayer match
			case packets.OsuMatchCreate:
				match := packets.ReadMultiplayerMatch(packetDataReader)

				lobby.CreateNewMultiMatch(match, client)
				break
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
				break
			//The client wants to change in which multiplayer slot its in
			case packets.OsuMatchChangeSlot:
				if client.currentMultiLobby != nil {
					var newSlot int32

					binary.Read(packetDataReader, binary.LittleEndian, &newSlot)

					client.currentMultiLobby.TryChangeSlot(client, int(newSlot))
				}
				break
			//The client wants to change sides
			case packets.OsuMatchChangeTeam:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.ChangeTeam(client)
				}
				break
			//The client wants to transfer its host status onto someone else
			case packets.OsuMatchTransferHost:
				if client.currentMultiLobby != nil {
					var newHost int32

					binary.Read(packetDataReader, binary.LittleEndian, &newHost)

					client.currentMultiLobby.TransferHost(client, int(newHost))
				}
				break
			//The client informs the server it has pressed the ready button
			case packets.OsuMatchReady:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.ReadyUp(client)
				}
				break
			//The client informs the server it has pressed the not ready button
			case packets.OsuMatchNotReady:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.Unready(client)
				}
				break
			//The client informs the server it has made some changes to the settings of the match
			case packets.OsuMatchChangeSettings:
				if client.currentMultiLobby != nil {
					newMatch := packets.ReadMultiplayerMatch(packetDataReader)
					client.currentMultiLobby.ChangeSettings(client, newMatch)
				}
				break
			//The client informs the server that it has changed the mods in the match
			case packets.OsuMatchChangeMods:
				if client.currentMultiLobby != nil {
					var newMods int32

					binary.Read(packetDataReader, binary.LittleEndian, &newMods)

					client.currentMultiLobby.ChangeMods(client, newMods)
				}
				break
			//The client informs the server that it has tried to lock a slot in the multi lobby
			case packets.OsuMatchLock:
				if client.currentMultiLobby != nil {
					var slot int32

					binary.Read(packetDataReader, binary.LittleEndian, &slot)

					client.currentMultiLobby.LockSlot(client, int(slot))
				}
				break
			//The client informs the server that it's missing the beatmap to be played
			case packets.OsuMatchNoBeatmap:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.InformNoBeatmap(client)
				}
				break
			//The client informs the server that it now has the beatmap that is to be played
			case packets.OsuMatchHasBeatmap:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.InformGotBeatmap(client)
				}
				break
			//The client informs the server that it has completed playing the map
			case packets.OsuMatchComplete:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.InformCompletion(client)
				}
				break
			//The client informs the server that it has loaded into the game successfully
			case packets.OsuMatchLoadComplete:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.InformLoadComplete(client)
				}
				break
			//The client informs the server of its new score
			case packets.OsuMatchScoreUpdate:
				if client.currentMultiLobby != nil {
					scoreFrame := packets.ReadScoreFrame(packetDataReader)

					client.currentMultiLobby.InformScoreUpdate(client, scoreFrame)
				}
				break
			//The client has requested to skip the beginning break
			case packets.OsuMatchSkipRequest:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.InformPressedSkip(client)
				}
				break
			//The client has failed the map
			case packets.OsuMatchFailed:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.InformFailed(client)
				}
				break
			//The client has pressed start game
			case packets.OsuMatchStart:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.StartGame(client)
				}
				break
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
				go database.AddFriend(client.UserData.UserID, uint64(friendId))
				break
			//The client is looking to remove a friend from their friends list
			case packets.OsuFriendsRemove:
				var friendId int32

				binary.Read(packetDataReader, binary.LittleEndian, &friendId)

				for index, value := range client.FriendsList {
					if value.User2 == uint64(friendId) {
						client.FriendsList = append(client.FriendsList[0:index], client.FriendsList[index+1:]...)
					}
				}

				go database.RemoveFriend(client.UserData.UserID, uint64(friendId))
				break
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
				break
			default:
				fmt.Printf("%s: %s: Got %s, of Size: %d\n", time.Now().Format(time.ANSIC), client.UserData.Username, packets.GetPacketName(packet.PacketId), packet.PacketSize)
				break
			}
		}
	}
}

// SendOutgoing is looping over the packet queue and waiting for new packets, and sends them off as they come in
func (client *Client) SendOutgoing() {
	for packet := range client.PacketQueue {
		if packet.PacketId != 8 {
			fmt.Printf("Sending %s to %s\n", packets.GetPacketName(packet.PacketId), client.UserData.Username)
		}

		client.connection.Write(packet.GetBytes())
	}
}

// MaintainClient is looping every second, sending out pings and handles timeouts
func (client *Client) MaintainClient() {
	for client.continueRunning {
		if client.lastReceive.Add(ReceiveTimeout).Before(time.Now()) {
			fmt.Printf("%s Timed out!\n", client.UserData.Username)
			fmt.Printf("time.Now(): %s, lastRecieve.Add(RecieveTimeout): %s; lastRecieve: %s", time.Now(), client.lastReceive.Add(ReceiveTimeout).String(), client.lastReceive.String())

			client.CleanupClient()

			client.continueRunning = false
		}

		if client.lastPing.Add(PingTimeout).Before(time.Now()) {
			packets.BanchoSendPing(client.PacketQueue)

			client.lastPing = time.Now()
		}

		time.Sleep(time.Second)
	}

	//We close in MaintainClient instead of in CleanupClient to avoid possible double closes, causing panics
	close(client.PacketQueue)
	client.connection.Close()
}
