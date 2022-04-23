package packets

func BanchoSendPing(packetQueue chan BanchoPacket) {
	BanchoSendEmptyPacket(packetQueue, BanchoPing)
}
