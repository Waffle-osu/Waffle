package packets

import (
	"bytes"
	"container/list"
	"encoding/binary"
)

func BanchoSendChannelJoinSuccess(packetQueue *list.List, channelName string) {
	buf := new(bytes.Buffer)

	//Write Data
	binary.Write(buf, binary.LittleEndian, WriteBanchoString(channelName))

	packetBytes := buf.Bytes()
	packetLength := len(packetBytes)

	packet := BanchoPacket{
		PacketId:          BanchoChannelJoinSuccess,
		PacketCompression: 0,
		PacketSize:        int32(packetLength),
		PacketData:        packetBytes,
	}

	packetQueue.PushBack(packet)
}
