package packets

func BanchoSendMatchAllPlayersLoaded(packetQueue chan BanchoPacket) {
	BanchoSendEmptyPacket(packetQueue, BanchoMatchAllPlayersLoaded)
}
