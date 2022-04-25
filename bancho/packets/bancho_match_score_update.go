package packets

import "bytes"

func BanchoSendMatchScoreUpdate(packetQueue chan BanchoPacket, scoreFrame ScoreFrame) {
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
