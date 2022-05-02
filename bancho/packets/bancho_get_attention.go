package packets

func BanchoSendGetAttention(packetQueue chan BanchoPacket) {
	BanchoSendEmptyPacket(packetQueue, BanchoGetAttention)
}
