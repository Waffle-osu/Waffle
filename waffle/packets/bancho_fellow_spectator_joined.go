package packets

func BanchoSendFellowSpectatorJoined(packetQueue chan BanchoPacket, userId int32) {
	BanchoSendIntPacket(packetQueue, BanchoFellowSpectatorJoined, userId)
}
