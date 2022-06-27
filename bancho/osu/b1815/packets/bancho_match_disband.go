package packets

func BanchoSendMatchDisband(packetQueue chan BanchoPacket, matchId int32) {
	BanchoSendIntPacket(packetQueue, BanchoMatchDisband, matchId)
}
