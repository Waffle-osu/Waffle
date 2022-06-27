package packets

import (
	"Waffle/bancho/osu/base_packet_structures"
	"bytes"
)

func BanchoSendBeatmapInfoReply(packetQueue chan BanchoPacket, infoReply base_packet_structures.BeatmapInfoReply) {
	buf := new(bytes.Buffer)

	infoReply.WriteBeatmapInfoReply(buf)

	packetBytes := buf.Bytes()
	packetLength := len(packetBytes)

	packet := BanchoPacket{
		PacketId:          BanchoBeatmapInfoReply,
		PacketCompression: 0,
		PacketSize:        int32(packetLength),
		PacketData:        packetBytes,
	}

	packetQueue <- packet
}
