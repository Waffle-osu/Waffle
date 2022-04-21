package clients

import (
	"Waffle/waffle/database"
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

			if packet.PacketId == 4 {
				continue
			}

			packetDataReader := bytes.NewBuffer(packet.PacketData)

			switch packet.PacketId {
			case packets.OsuRequestStatusUpdate:
				var stats database.UserStats

				switch client.Status.CurrentPlaymode {
				case packets.OsuGamemodeOsu:
					stats = client.OsuStats
					break
				case packets.OsuGamemodeTaiko:
					stats = client.TaikoStats
					break
				case packets.OsuGamemodeCatch:
					stats = client.CatchStats
					break
				case packets.OsuGamemodeMania:
					stats = client.ManiaStats
					break
				}

				packets.BanchoSendUserPresence(client.PacketQueue, client.UserData, client.OsuStats, client.ClientData.Timezone)
				packets.BanchoSendOsuUpdate(client.PacketQueue, stats, client.Status)
				break
			case packets.OsuUserStatsRequest:
				var listLength int16

				binary.Read(packetDataReader, binary.LittleEndian, &listLength)

				for i := 0; int16(i) != listLength; i++ {
					var currentId int32
					binary.Read(packetDataReader, binary.LittleEndian, currentId)

					user := GetClientById(currentId)

					if user == nil {
						continue
					}

					var stats database.UserStats

					switch user.Status.CurrentPlaymode {
					case packets.OsuGamemodeOsu:
						stats = user.OsuStats
						break
					case packets.OsuGamemodeTaiko:
						stats = user.TaikoStats
						break
					case packets.OsuGamemodeCatch:
						stats = user.CatchStats
						break
					case packets.OsuGamemodeMania:
						stats = user.ManiaStats
						break
					}

					packets.BanchoSendOsuUpdate(client.PacketQueue, stats, user.Status)
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

	close(client.PacketQueue)
}
