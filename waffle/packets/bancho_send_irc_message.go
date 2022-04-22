package packets

import (
	"bytes"
	"encoding/binary"
)

func BanchoSendIrcMessage(packetQueue chan BanchoPacket, senderUsername string, target string, content string) {
	buf := new(bytes.Buffer)

	//Write Data
	binary.Write(buf, binary.LittleEndian, WriteBanchoString(senderUsername))
	binary.Write(buf, binary.LittleEndian, WriteBanchoString(content))
	binary.Write(buf, binary.LittleEndian, WriteBanchoString(target))

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
