package packets

import (
	"bytes"
	"encoding/binary"
)

func BanchoSendHandleOsuQuit(packetQueue chan BanchoPacket, userId int32) {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, userId)

	packetBytes := buf.Bytes()
	packetLength := len(packetBytes)

	packet := BanchoPacket{
		PacketId:          BanchoHandleOsuQuit,
		PacketCompression: 0,
		PacketSize:        int32(packetLength),
		PacketData:        packetBytes,
	}

	packetQueue <- packet
}
