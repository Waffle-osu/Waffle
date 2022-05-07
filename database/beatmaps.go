package database

type Beatmap struct {
	BeatmapId     int32
	BeatmapsetId  int32
	CreatorId     int64
	Filename      string
	BeatmapMd5    string
	Version       string
	TotalLength   int32
	DrainTime     int32
	CountObjects  int32
	CountNormal   int32
	CountSlider   int32
	CountSpinner  int32
	DiffHp        int8
	DiffCs        int8
	DiffOd        int8
	DiffStars     float32
	Playmode      int8
	RankingStatus int8
	LastUpdate    string
	SubmitDate    string
	ApproveDate   string
	BeatmapSource int8
}
