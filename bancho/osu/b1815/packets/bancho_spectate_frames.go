package packets

import (
	"Waffle/bancho/osu/base_packet_structures"
	"bytes"
)

func BanchoSendSpectateFrames(packetQueue chan BanchoPacket, frameBundle base_packet_structures.SpectatorFrameBundle) {
	buf := new(bytes.Buffer)

	frameBundle.WriteSpectatorFrameBundle(buf)

	packetBytes := buf.Bytes()
	packetLength := len(packetBytes)

	packet := BanchoPacket{
		PacketId:          BanchoSpectateFrames,
		PacketCompression: 0,
		PacketSize:        int32(packetLength),
		PacketData:        packetBytes,
	}

	packetQueue <- packet
}
