package packets

func BanchoSendPing(packetQueue chan BanchoPacket) {
	packet := BanchoPacket{
		PacketId:          BanchoPing,
		PacketCompression: 0,
		PacketSize:        0,
		PacketData:        []byte{},
	}

	packetQueue <- packet
}
