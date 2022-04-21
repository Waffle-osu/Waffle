package packets

import (
	"bytes"
	"container/list"
)

func BasePacket(packetQueue *list.List) {
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

	packetQueue.PushBack(packet)
}
