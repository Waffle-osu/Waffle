package base_packet_structures

type BeatmapInfo struct {
	InfoId          int16
	BeatmapId       int32
	BeatmapSetId    int32
	ThreadId        int32
	Ranked          uint8
	OsuRank         uint8
	TaikoRank       uint8
	CatchRank       uint8
	BeatmapChecksum string
}
