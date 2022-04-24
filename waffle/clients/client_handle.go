package clients

import (
	"Waffle/waffle/chat"
	"Waffle/waffle/client_manager"
	"Waffle/waffle/database"
	"Waffle/waffle/lobby"
	"Waffle/waffle/packets"
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

func (client *Client) HandleIncoming() {
	readBuffer := make([]byte, 4096)

	for client.continueRunning {
		read, readErr := client.connection.Read(readBuffer)

		if readErr != nil {
			return
		}

		client.lastReceive = time.Now()

		//Get the bytes that were actually read
		packetBuffer := bytes.NewBuffer(readBuffer[:read])
		readIndex := 0

		for readIndex < read {
			read, packet := packets.ReadBanchoPacketHeader(packetBuffer)

			readIndex += read

			//Unused packet
			if packet.PacketId == 79 {
				continue
			}

			packetDataReader := bytes.NewBuffer(packet.PacketData)

			switch packet.PacketId {
			case packets.OsuSendUserStatus:
				statusUpdate := packets.ReadStatusUpdate(packetDataReader)

				client.Status = statusUpdate

				client_manager.BroadcastPacket(func(packetQueue chan packets.BanchoPacket) {
					packets.BanchoSendOsuUpdate(packetQueue, client.OsuStats, client.Status)
				})
				break
			case packets.OsuRequestStatusUpdate:
				packets.BanchoSendUserPresence(client.PacketQueue, client.UserData, client.OsuStats, client.GetClientTimezone())
				packets.BanchoSendOsuUpdate(client.PacketQueue, client.GetRelevantUserStats(), client.Status)
				break
			case packets.OsuUserStatsRequest:
				var listLength int16

				binary.Read(packetDataReader, binary.LittleEndian, &listLength)

				for i := 0; int16(i) != listLength; i++ {
					var currentId int32
					binary.Read(packetDataReader, binary.LittleEndian, currentId)

					user := client_manager.GetClientById(currentId)

					if user == nil {
						continue
					}

					packets.BanchoSendOsuUpdate(client.PacketQueue, user.GetRelevantUserStats(), user.GetUserStatus())
					break
				}
			case packets.OsuSendIrcMessage:
				message := packets.ReadMessage(packetDataReader)

				if message.Target == "#multiplayer" {
					client.currentMultiLobby.MultiChannel.SendMessage(client, message.Message, message.Target)
					break
				}

				for _, channel := range client.joinedChannels {
					if channel.Name == message.Target {
						channel.SendMessage(client, message.Message, message.Target)
					}
				}

				break
			case packets.OsuSendIrcMessagePrivate:
				message := packets.ReadMessage(packetDataReader)
				message.Sender = client.UserData.Username

				targetClient := client_manager.GetClientByName(message.Target)

				if targetClient != nil {
					packets.BanchoSendIrcMessage(targetClient.GetPacketQueue(), message)

					awayMessage := targetClient.GetAwayMessage()

					if awayMessage != "" {
						packets.BanchoSendIrcMessage(client.PacketQueue, packets.Message{
							Sender:  targetClient.GetUserData().Username,
							Message: fmt.Sprintf("/me is away! %s", awayMessage),
							Target:  client.GetUserData().Username,
						})
					}
				}
				break
			case packets.OsuExit:
				client.CleanupClient()
				return
			case packets.OsuStartSpectating:
				var spectatorId int32

				binary.Read(packetDataReader, binary.LittleEndian, &spectatorId)

				toSpectate := client_manager.GetClientById(spectatorId)

				if toSpectate == nil {
					break
				}

				toSpectate.InformSpectatorJoin(client)

				client.spectatingClient = toSpectate
				break
			case packets.OsuStopSpectating:
				if client.spectatingClient == nil {
					break
				}

				client.spectatingClient.InformSpectatorLeft(client)
				client.spectatingClient = nil
				break
			case packets.OsuSpectateFrames:
				frameBundle := packets.ReadSpectatorFrameBundle(packetDataReader)

				client.BroadcastToSpectators(func(packetQueue chan packets.BanchoPacket) {
					packets.BanchoSendSpectateFrames(packetQueue, frameBundle)
				})
				break
			case packets.OsuCantSpectate:
				if client.spectatingClient != nil {
					client.spectatingClient.InformSpectatorCantSpectate(client)
				}
				break
			case packets.OsuErrorReport:
				errorString := string(packets.ReadBanchoString(packetDataReader))

				fmt.Printf("%s Encountered an error!! Error Details:\n%s", client.UserData.Username, errorString)
				break
			case packets.OsuPong:
				client.lastReceive = time.Now()
				break
			case packets.OsuLobbyJoin:
				lobby.JoinLobby(client)
				client.isInLobby = true
				break
			case packets.OsuLobbyPart:
				lobby.PartLobby(client)
				client.isInLobby = false
				break
			case packets.OsuChannelJoin:
				channelName := string(packets.ReadBanchoString(packetDataReader))

				channel, exists := chat.GetChannelByName(channelName)

				if exists {
					if channel.Join(client) {
						packets.BanchoSendChannelJoinSuccess(client.PacketQueue, channelName)
						client.joinedChannels = append(client.joinedChannels, channel)
					} else {
						packets.BanchoSendChannelRevoked(client.PacketQueue, channelName)
					}
				} else {
					packets.BanchoSendChannelRevoked(client.PacketQueue, channelName)
				}
				break
			case packets.OsuChannelLeave:
				channelName := string(packets.ReadBanchoString(packetDataReader))

				for index, channel := range client.joinedChannels {
					if channel.Name == channelName {
						channel.Leave(client)
						client.joinedChannels = append(client.joinedChannels[0:index], client.joinedChannels[index+1:]...)
					}
				}
				break
			case packets.OsuMatchCreate:
				match := packets.ReadMultiplayerMatch(packetDataReader)

				lobby.CreateNewMultiMatch(match, client)
				break
			case packets.OsuMatchJoin:
				matchJoin := packets.ReadMatchJoin(packetDataReader)

				foundMatch := lobby.GetMultiMatchById(uint16(matchJoin.MatchId))

				if foundMatch != nil {
					client.JoinMatch(foundMatch, matchJoin.Password)
				} else {
					packets.BanchoSendMatchJoinFail(client.PacketQueue)
				}
			case packets.OsuMatchPart:
				client.LeaveCurrentMatch()
				break
			case packets.OsuMatchChangeSlot:
				if client.currentMultiLobby != nil {
					var newSlot int32

					binary.Read(packetDataReader, binary.LittleEndian, &newSlot)

					client.currentMultiLobby.TryChangeSlot(client, int(newSlot))
				}
				break
			case packets.OsuMatchChangeTeam:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.ChangeTeam(client)
				}
				break
			case packets.OsuMatchTransferHost:
				if client.currentMultiLobby != nil {
					var newHost int32

					binary.Read(packetDataReader, binary.LittleEndian, &newHost)

					client.currentMultiLobby.TransferHost(client, int(newHost))
				}
				break
			case packets.OsuMatchReady:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.ReadyUp(client)
				}
				break
			case packets.OsuMatchNotReady:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.Unready(client)
				}
				break
			case packets.OsuMatchChangeSettings:
				if client.currentMultiLobby != nil {
					newMatch := packets.ReadMultiplayerMatch(packetDataReader)
					client.currentMultiLobby.ChangeSettings(client, newMatch)
				}
				break
			case packets.OsuMatchChangeMods:
				if client.currentMultiLobby != nil {
					var newMods int32

					binary.Read(packetDataReader, binary.LittleEndian, &newMods)

					client.currentMultiLobby.ChangeMods(client, newMods)
				}
				break
			case packets.OsuMatchLock:
				if client.currentMultiLobby != nil {
					var slot int32

					binary.Read(packetDataReader, binary.LittleEndian, &slot)

					client.currentMultiLobby.LockSlot(client, int(slot))
				}
				break
			case packets.OsuMatchNoBeatmap:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.InformNoBeatmap(client)
				}
				break
			case packets.OsuMatchHasBeatmap:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.InformGotBeatmap(client)
				}
				break
			case packets.OsuMatchComplete:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.InformCompletion(client)
				}
				break
			case packets.OsuMatchLoadComplete:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.InformLoadComplete(client)
				}
				break
			case packets.OsuMatchScoreUpdate:
				if client.currentMultiLobby != nil {
					scoreFrame := packets.ReadScoreFrame(packetDataReader)

					client.currentMultiLobby.InformScoreUpdate(client, scoreFrame)
				}
				break
			case packets.OsuMatchSkipRequest:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.InformPressedSkip(client)
				}
				break
			case packets.OsuMatchFailed:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.InformFailed(client)
				}
				break
			case packets.OsuMatchStart:
				if client.currentMultiLobby != nil {
					client.currentMultiLobby.StartGame(client)
				}
				break
			case packets.OsuFriendsAdd:
				var friendId int32

				binary.Read(packetDataReader, binary.LittleEndian, &friendId)

				client.FriendsList = append(client.FriendsList, database.FriendEntry{
					User1: client.UserData.UserID,
					User2: uint64(friendId),
				})

				go database.AddFriend(client.UserData.UserID, uint64(friendId))
				break
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
			case packets.OsuSetIrcAwayMessage:
				awayMessage := packets.ReadMessage(packetDataReader)

				client.awayMessage = awayMessage.Message

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
			}

			fmt.Printf("%s: Got %s, of Size: %d\n", client.UserData.Username, packets.GetPacketName(packet.PacketId), packet.PacketSize)
		}
	}
}

func (client *Client) SendOutgoing() {
	for packet := range client.PacketQueue {
		if packet.PacketId != 8 {
			fmt.Printf("Sending %s to %s\n", packets.GetPacketName(packet.PacketId), client.UserData.Username)
		}

		client.connection.Write(packet.GetBytes())
	}
}

func (client *Client) MaintainClient() {
	for client.continueRunning {
		if client.lastReceive.Add(ReceiveTimeout).Before(time.Now()) {
			fmt.Printf("%s Timed out!\n", client.UserData.Username)

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
