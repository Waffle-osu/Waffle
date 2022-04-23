package packets

func BanchoSendSpectatorLeft(packetQueue chan BanchoPacket, userId int32) {
	BanchoSendIntPacket(packetQueue, BanchoSpectatorLeft, userId)
}
