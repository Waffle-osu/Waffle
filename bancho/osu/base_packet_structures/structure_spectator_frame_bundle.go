package base_packet_structures

import (
	"encoding/binary"
	"io"
)

type SpectatorFrameBundle struct {
	FrameCount   uint16
	Frames       []SpectatorFrame `length:"FrameCount"`
	ReplayAction uint8
	ScoreFrame   ScoreFrame
}

func ReadSpectatorFrameBundle(reader io.Reader) SpectatorFrameBundle {
	frameBundle := SpectatorFrameBundle{}

	binary.Read(reader, binary.LittleEndian, &frameBundle.FrameCount)

	for i := 0; i != int(frameBundle.FrameCount); i++ {
		frameBundle.Frames = append(frameBundle.Frames, ReadSpectatorFrame(reader))
	}

	binary.Read(reader, binary.LittleEndian, &frameBundle.ReplayAction)

	frameBundle.ScoreFrame = ReadScoreFrame(reader)

	return frameBundle
}

func (frameBundle SpectatorFrameBundle) Write(writer io.Writer) {
	binary.Write(writer, binary.LittleEndian, frameBundle.FrameCount)

	for i := 0; i != int(frameBundle.FrameCount); i++ {
		frameBundle.Frames[i].Write(writer)
	}

	binary.Write(writer, binary.LittleEndian, frameBundle.ReplayAction)

	frameBundle.ScoreFrame.Write(writer)
}
