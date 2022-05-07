package database

type Beatmapset struct {
	BeatmapsetId  int32
	CreatorId     int64
	Artist        string
	Title         string
	Creator       string
	Source        string
	Tags          string
	HasVideo      int8
	HasStoryboard int8
	Bpm           float32
}
