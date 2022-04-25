package packets

func BanchoSendLobbyJoin(packetQueue chan BanchoPacket, userId int32) {
	BanchoSendIntPacket(packetQueue, BanchoLobbyJoin, userId)
}
