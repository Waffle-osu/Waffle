package clients

import (
	"Waffle/waffle/client_manager"
	"Waffle/waffle/packet_structures"
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
				statusUpdate := packet_structures.ReadStatusUpdate(packetDataReader)

				client.Status.CurrentStatus = statusUpdate.Status
				client.Status.StatusText = statusUpdate.StatusText
				client.Status.BeatmapChecksum = statusUpdate.BeatmapChecksum
				client.Status.CurrentMods = statusUpdate.CurrentMods
				client.Status.CurrentPlaymode = statusUpdate.Playmode
				client.Status.BeatmapId = statusUpdate.BeatmapId

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
				var message string
				var target string

				packets.ReadBanchoString(packetDataReader) //We get the Username from the client, no need for this though the client sends it anyway so we gotta read
				message = string(packets.ReadBanchoString(packetDataReader))
				target = string(packets.ReadBanchoString(packetDataReader))

				for _, channel := range client.joinedChannels {
					if channel.Name == target {
						channel.SendMessage(client, message, target)
					}
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
				frameBundle := packet_structures.ReadSpectatorFrameBundle(packetDataReader)

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
			default:
				fmt.Printf("Read Packet ID: %d, of Size: %d, current readIndex: %d\n", packet.PacketId, packet.PacketSize, readIndex)
			}
		}
	}
}

func (client *Client) SendOutgoing() {
	for packet := range client.PacketQueue {
		if packet.PacketId != 8 {
			fmt.Printf("Sending Packet %d to %s\n", packet.PacketId, client.UserData.Username)
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
