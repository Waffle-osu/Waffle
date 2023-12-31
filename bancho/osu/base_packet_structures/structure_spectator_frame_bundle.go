package base_packet_structures

type SpectatorFrameBundle struct {
	FrameCount   uint16
	Frames       []SpectatorFrame `length:"FrameCount"`
	ReplayAction uint8
	ScoreFrame   ScoreFrame
}
