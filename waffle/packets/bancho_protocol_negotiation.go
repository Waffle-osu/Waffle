package packets

import (
	"bytes"
	"container/list"
	"encoding/binary"
)

const Build1816ProtocolVersion int32 = 7

func BanchoSendProtocolNegotiation(packetQueue *list.List) {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, Build1816ProtocolVersion)

	packetBytes := buf.Bytes()
	packetLength := len(packetBytes)

	packet := BanchoPacket{
		PacketId:          BanchoProtocolNegotiation,
		PacketCompression: 0,
		PacketSize:        int32(packetLength),
		PacketData:        packetBytes,
	}

	packetQueue.PushBack(packet)
}
