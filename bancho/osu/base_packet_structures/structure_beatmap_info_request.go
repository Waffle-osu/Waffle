package base_packet_structures

type BeatmapInfoRequest struct {
	FilenameCount int32
	Filenames     []string `length:"FilenameCount"`
	IdCount       int32
	BeatmapIds    []int32 `length:"IdCount"`
}
