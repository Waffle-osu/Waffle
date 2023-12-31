package base_packet_structures

type BeatmapInfoReply struct {
	Count        int32
	BeatmapInfos []BeatmapInfo `length:"Count"`
}
