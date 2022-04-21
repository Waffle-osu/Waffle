package clients

import (
	"Waffle/waffle/packets"
	"bytes"
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

			fmt.Printf("Read Packet ID: %d, of Size: %d, current readIndex: %d\n", packet.PacketId, packet.PacketSize, readIndex)

			//switch packet.PacketId {
			//
			//}
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
}
