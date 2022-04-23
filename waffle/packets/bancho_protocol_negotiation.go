package packets

const Build1816ProtocolVersion int32 = 7

func BanchoSendProtocolNegotiation(packetQueue chan BanchoPacket) {
	BanchoSendIntPacket(packetQueue, BanchoProtocolNegotiation, Build1816ProtocolVersion)
}
