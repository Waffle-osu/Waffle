package packets

import (
	"Waffle/helpers/serialization"
	"bytes"
	"encoding/binary"
)

func BanchoSendChannelJoinSuccess(packetQueue chan BanchoPacket, channelName string) {
	buf := new(bytes.Buffer)

	//Write Data
	binary.Write(buf, binary.LittleEndian, serialization.WriteBanchoString(channelName))

	packetBytes := buf.Bytes()
	packetLength := len(packetBytes)

	packet := BanchoPacket{
		PacketId:          BanchoChannelJoinSuccess,
		PacketCompression: 0,
		PacketSize:        int32(packetLength),
		PacketData:        packetBytes,
	}

	packetQueue <- packet
}
