package packets

func BanchoSendMatchComplete(packetQueue chan BanchoPacket) {
	BanchoSendEmptyPacket(packetQueue, BanchoMatchComplete)
}
