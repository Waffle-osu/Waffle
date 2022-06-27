package packets

import (
	"Waffle/bancho/osu/base_packet_structures"
	"bytes"
)

func BanchoSendMatchUpdate(packetQueue chan BanchoPacket, match base_packet_structures.MultiplayerMatch) {
	buf := new(bytes.Buffer)

	match.WriteMultiplayerMatch(buf)

	packetBytes := buf.Bytes()
	packetLength := len(packetBytes)

	packet := BanchoPacket{
		PacketId:          BanchoMatchUpdate,
		PacketCompression: 0,
		PacketSize:        int32(packetLength),
		PacketData:        packetBytes,
	}

	packetQueue <- packet
}
