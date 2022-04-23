package packets

func BanchoSendMatchJoinFail(packetQueue chan BanchoPacket) {
	BanchoSendIntPacket(packetQueue, BanchoMatchJoinFail, 0)
}
