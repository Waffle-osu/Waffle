package clients

import (
	"Waffle/waffle/packets"
	"bytes"
	"fmt"
)

func (client *Client) HandleIncoming() {
	readBuffer := make([]byte, 4096)

	//Check if there's at least 1 packet header there
	availableBytes := client.BufReader.Buffered()

	if availableBytes > 0 {
		read, readErr := client.Connection.Read(readBuffer)

		if readErr != nil {
			return
		}

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
	sendBuffer := new(bytes.Buffer)

	client.PacketQueueMutex.Lock()

	for retrievedPacket := client.PacketQueue.Front(); retrievedPacket != nil; retrievedPacket = retrievedPacket.Next() {
		packet := retrievedPacket.Value.(packets.BanchoPacket)

		fmt.Printf("Sending Packet %d\n", packet.PacketId)

		sendBuffer.Write(packet.GetBytes())

		client.PacketQueue.Remove(retrievedPacket)
	}

	client.PacketQueueMutex.Unlock()

	client.Connection.Write(sendBuffer.Bytes())
}
