package clients

import (
	"Waffle/waffle/client_manager"
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

			if packet.PacketId == 4 || packet.PacketId == 79 {
				continue
			}

			packetDataReader := bytes.NewBuffer(packet.PacketData)

			switch packet.PacketId {
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
