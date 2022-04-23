package packets

func BanchoSendFellowSpectatorLeft(packetQueue chan BanchoPacket, userId int32) {
	BanchoSendIntPacket(packetQueue, BanchoFellowSpectatorLeft, userId)
}
