package packets

import "bytes"

func BanchoSendMatchJoinSuccess(packetQueue chan BanchoPacket, match MultiplayerMatch) {
	buf := new(bytes.Buffer)

	match.WriteMultiplayerMatch(buf)

	packetBytes := buf.Bytes()
	packetLength := len(packetBytes)

	packet := BanchoPacket{
		PacketId:          BanchoMatchJoinSuccess,
		PacketCompression: 0,
		PacketSize:        int32(packetLength),
		PacketData:        packetBytes,
	}

	packetQueue <- packet
}
