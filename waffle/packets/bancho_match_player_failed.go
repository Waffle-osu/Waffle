package packets

func BanchoSendMatchPlayerFailed(packetQueue chan BanchoPacket, slot int32) {
	BanchoSendIntPacket(packetQueue, BanchoMatchPlayerFailed, slot)
}
