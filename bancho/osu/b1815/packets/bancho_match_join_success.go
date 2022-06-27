package packets

import (
	"Waffle/bancho/osu/base_packet_structures"
	"bytes"
)

func BanchoSendMatchJoinSuccess(packetQueue chan BanchoPacket, match base_packet_structures.MultiplayerMatch) {
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
