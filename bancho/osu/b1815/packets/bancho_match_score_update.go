package packets

import (
	"Waffle/bancho/osu/base_packet_structures"
	"bytes"
)

func BanchoSendMatchScoreUpdate(packetQueue chan BanchoPacket, scoreFrame base_packet_structures.ScoreFrame) {
	buf := new(bytes.Buffer)

	scoreFrame.WriteScoreFrame(buf)

	packetBytes := buf.Bytes()
	packetLength := len(packetBytes)

	packet := BanchoPacket{
		PacketId:          BanchoMatchScoreUpdate,
		PacketCompression: 0,
		PacketSize:        int32(packetLength),
		PacketData:        packetBytes,
	}

	packetQueue <- packet
}
