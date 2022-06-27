package packets

func BanchoSendSpectatorJoin(packetQueue chan BanchoPacket, userId int32) {
	BanchoSendIntPacket(packetQueue, BanchoSpectatorJoined, userId)
}
