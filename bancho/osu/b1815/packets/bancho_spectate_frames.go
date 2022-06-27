package packets

import (
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/helpers/serialization"
	"bytes"
)

func BanchoSendSpectateFrames(packetQueue chan serialization.BanchoPacket, frameBundle base_packet_structures.SpectatorFrameBundle) {
	buf := new(bytes.Buffer)

	frameBundle.WriteSpectatorFrameBundle(buf)

	packetBytes := buf.Bytes()
	packetLength := len(packetBytes)

	packet := serialization.BanchoPacket{
		PacketId:          serialization.BanchoSpectateFrames,
		PacketCompression: 0,
		PacketSize:        int32(packetLength),
		PacketData:        packetBytes,
	}

	packetQueue <- packet
}
