package packets

func BanchoSendMatchPlayerSkipped(packetQueue chan BanchoPacket, slot int32) {
	BanchoSendIntPacket(packetQueue, BanchoMatchPlayerSkipped, slot)
}
