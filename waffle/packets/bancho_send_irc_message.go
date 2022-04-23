package packets

import (
	"bytes"
)

func BanchoSendIrcMessage(packetQueue chan BanchoPacket, message Message) {
	buf := new(bytes.Buffer)

	message.WriteMessage(buf)

	packetBytes := buf.Bytes()
	packetLength := len(packetBytes)

	packet := BanchoPacket{
		PacketId:          BanchoSendMessage,
		PacketCompression: 0,
		PacketSize:        int32(packetLength),
		PacketData:        packetBytes,
	}

	packetQueue <- packet
}
