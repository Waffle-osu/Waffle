package packets

func BanchoSendSpectatorCantSpectate(packetQueue chan BanchoPacket, userId int32) {
	BanchoSendIntPacket(packetQueue, BanchoSpectatorCantSpectate, userId)
}
