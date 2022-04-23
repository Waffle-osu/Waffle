package clients

import (
	"Waffle/waffle/chat"
	"Waffle/waffle/client_manager"
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
			CleanupClient(client)
			return
		}

		client.lastReceive = time.Now()

		//Get the bytes that were actually read
		packetBuffer := bytes.NewBuffer(readBuffer[:read])
		readIndex := 0

		for readIndex < read {
			read, packet := packets.ReadBanchoPacketHeader(packetBuffer)

			readIndex += read

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

				for _, channel := range client.joinedChannels {
					if channel.Name == message.Target {
						channel.SendMessage(client, message.Message, message.Target)
					}
				}

				break
			case packets.OsuSendIrcMessagePrivate:
				message := packets.ReadMessage(packetDataReader)

				targetClient := client_manager.GetClientByName(message.Target)

				if targetClient != nil {
					packets.BanchoSendIrcMessage(targetClient.GetPacketQueue(), message)
				}
				break
			case packets.OsuExit:
				CleanupClient(client)
				break
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
			case packets.OsuMatchPart:
				client.LeaveCurrentMatch()
				break
			default:
				fmt.Printf("Got %s, of Size: %d\n", packets.GetPacketName(packet.PacketId), packet.PacketSize)
			}
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
			CleanupClient(client)
		}

		if client.lastPing.Add(PingTimeout).Before(time.Now()) {
			packets.BanchoSendPing(client.PacketQueue)

			client.lastPing = time.Now()
		}

		time.Sleep(time.Second)
	}

	//We close in MaintainClient instead of in CleanupClient to avoid possible double closes, causing panics
	close(client.PacketQueue)
}
