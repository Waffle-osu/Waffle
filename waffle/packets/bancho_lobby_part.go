package packets

func BanchoSendLobbyPart(packetQueue chan BanchoPacket, userId int32) {
	BanchoSendIntPacket(packetQueue, BanchoLobbyPart, userId)
}
