package packets

func BanchoSendMatchSkip(packetQueue chan BanchoPacket) {
	BanchoSendEmptyPacket(packetQueue, BanchoMatchSkip)
}
