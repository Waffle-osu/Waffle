package packets

import (
	"bytes"
)

func BanchoSendBeatmapInfoReply(packetQueue chan BanchoPacket, infoReply BeatmapInfoReply) {
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
