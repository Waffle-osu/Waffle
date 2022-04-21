package packets

import (
	"Waffle/waffle/chat"
	"bytes"
	"encoding/binary"
)

func BanchoSendChannelAvailable(packetQueue chan BanchoPacket, channel *chat.Channel) {
	buf := new(bytes.Buffer)

	//Write Data
	binary.Write(buf, binary.LittleEndian, WriteBanchoString(channel.Name))

	packetBytes := buf.Bytes()
	packetLength := len(packetBytes)

	packet := BanchoPacket{
		PacketId:          BanchoChannelAvailable,
		PacketCompression: 0,
		PacketSize:        int32(packetLength),
		PacketData:        packetBytes,
	}

	packetQueue <- packet
}
