package packets

import (
	"bytes"
)

func BanchoSendSpectateFrames(packetQueue chan BanchoPacket, frameBundle SpectatorFrameBundle) {
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
