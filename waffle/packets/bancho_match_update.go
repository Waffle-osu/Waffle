package packets

import "bytes"

func BanchoSendMatchUpdate(packetQueue chan BanchoPacket, match MultiplayerMatch) {
	buf := new(bytes.Buffer)

	match.WriteMultiplayerMatch(buf)

	packetBytes := buf.Bytes()
	packetLength := len(packetBytes)

	packet := BanchoPacket{
		PacketId:          BanchoMatchUpdate, //TODO: change this out
		PacketCompression: 0,
		PacketSize:        int32(packetLength),
		PacketData:        packetBytes,
	}

	packetQueue <- packet
}
