package packets

func BanchoSendMatchTransferHost(packetQueue chan BanchoPacket) {
	BanchoSendEmptyPacket(packetQueue, BanchoMatchTransferHost)
}
