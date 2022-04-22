package packets

import (
	"bytes"
	"encoding/binary"
)

func BanchoSendAnnounce(packetQueue chan BanchoPacket, announcement string) {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, WriteBanchoString(announcement))

	packetBytes := buf.Bytes()
	packetLength := len(packetBytes)

	packet := BanchoPacket{
		PacketId:          BanchoAnnounce,
		PacketCompression: 0,
		PacketSize:        int32(packetLength),
		PacketData:        packetBytes,
	}

	packetQueue <- packet
}
