package packets

import (
	"Waffle/bancho/osu/base_packet_structures"
	"bytes"
)

func BanchoSendMatchNew(packetQueue chan BanchoPacket, match base_packet_structures.MultiplayerMatch) {
	buf := new(bytes.Buffer)

	match.WriteMultiplayerMatch(buf)

	packetBytes := buf.Bytes()
	packetLength := len(packetBytes)

	packet := BanchoPacket{
		PacketId:          BanchoMatchNew,
		PacketCompression: 0,
		PacketSize:        int32(packetLength),
		PacketData:        packetBytes,
	}

	packetQueue <- packet
}
