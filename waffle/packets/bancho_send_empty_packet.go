package packets

import "bytes"

func BanchoSendEmptyPacket(packetQueue chan BanchoPacket, packetId uint16) {
	buf := new(bytes.Buffer)

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
