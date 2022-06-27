package packets

func BanchoSendHandleOsuQuit(packetQueue chan BanchoPacket, userId int32) {
	BanchoSendIntPacket(packetQueue, BanchoHandleOsuQuit, userId)
}
