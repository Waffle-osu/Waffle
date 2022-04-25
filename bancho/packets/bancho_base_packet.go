package packets

import (
	"bytes"
)

func BasePacket(packetQueue chan BanchoPacket) {
	buf := new(bytes.Buffer)

	//Write Data

	packetBytes := buf.Bytes()
	packetLength := len(packetBytes)

	packet := BanchoPacket{
		PacketId:          BanchoPing, //TODO: change this out
		PacketCompression: 0,
		PacketSize:        int32(packetLength),
		PacketData:        packetBytes,
	}

	packetQueue <- packet
}
