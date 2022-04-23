package packets

import (
	"bytes"
	"encoding/binary"
)

func BanchoSendIntPacket(packetQueue chan BanchoPacket, packetId uint16, integer int32) {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, integer)

	packetBytes := buf.Bytes()
	packetLength := len(packetBytes)

	packet := BanchoPacket{
		PacketId:          packetId,
		PacketCompression: 0,
		PacketSize:        int32(packetLength),
		PacketData:        packetBytes,
	}

	packetQueue <- packet
}
