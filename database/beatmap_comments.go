package database

type BeatmapComment struct {
	CommentId    int64
	UserId       uint64
	BeatmapId    int
	BeatmapSetId int
	ScoreId      uint64
	Time         int64
	Target       int8
	Comment      string
	FormatString string
}
