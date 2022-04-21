package packets

import (
	"Waffle/waffle/chat"
	"bytes"
	"container/list"
	"encoding/binary"
)

func BanchoSendChannelAvailable(packetQueue *list.List, channel *chat.Channel) {
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

	packetQueue.PushBack(packet)
}
