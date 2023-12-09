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

const (
	BeatmapsDatabaseStatusUnsubmitted = -1
	BeatmapsDatabaseStatusPending     = 0
	BeatmapsDatabaseStatusRanked      = 1
	BeatmapsDatabaseStatusApproved    = 2
)

func BeatmapsGetByMd5(checksum string) (queryResult int8, beatmap Beatmap) {
	beatmapQuery, beatmapQueryErr := Database.Query("SELECT beatmap_id, beatmapset_id, creator_id, filename, beatmap_md5, version, total_length, drain_time, count_objects, count_normal, count_slider, count_spinner, diff_hp, diff_cs, diff_od, diff_stars, playmode, ranking_status, last_update, submit_date, approve_date, beatmap_source FROM waffle.beatmaps WHERE beatmap_md5 = ?", checksum)

	if beatmapQueryErr != nil {
		if beatmapQuery != nil {
			beatmapQuery.Close()
		}

		return -2, Beatmap{}
	}

	if beatmapQuery.Next() {
		returnBeatmap := Beatmap{}

		scanErr := beatmapQuery.Scan(&returnBeatmap.BeatmapId, &returnBeatmap.BeatmapsetId, &returnBeatmap.CreatorId, &returnBeatmap.Filename, &returnBeatmap.BeatmapMd5, &returnBeatmap.Version, &returnBeatmap.TotalLength, &returnBeatmap.DrainTime, &returnBeatmap.CountObjects, &returnBeatmap.CountNormal, &returnBeatmap.CountSlider, &returnBeatmap.CountSpinner, &returnBeatmap.DiffHp, &returnBeatmap.DiffCs, &returnBeatmap.DiffOd, &returnBeatmap.DiffStars, &returnBeatmap.Playmode, &returnBeatmap.RankingStatus, &returnBeatmap.LastUpdate, &returnBeatmap.SubmitDate, &returnBeatmap.ApproveDate, &returnBeatmap.BeatmapSource)

		beatmapQuery.Close()

		if scanErr != nil {
			return -2, Beatmap{}
		}

		return 0, returnBeatmap
	} else {
		beatmapQuery.Close()
		return -1, Beatmap{}
	}
}

func BeatmapsGetByFilename(filename string) (queryResult int8, beatmap Beatmap) {
	beatmapQuery, beatmapQueryErr := Database.Query("SELECT beatmap_id, beatmapset_id, creator_id, filename, beatmap_md5, version, total_length, drain_time, count_objects, count_normal, count_slider, count_spinner, diff_hp, diff_cs, diff_od, diff_stars, playmode, ranking_status, last_update, submit_date, approve_date, beatmap_source FROM waffle.beatmaps WHERE filename = ?", filename)

	if beatmapQueryErr != nil {
		if beatmapQuery != nil {
			beatmapQuery.Close()
		}

		return -2, Beatmap{}
	}

	if beatmapQuery.Next() {
		returnBeatmap := Beatmap{}

		scanErr := beatmapQuery.Scan(&returnBeatmap.BeatmapId, &returnBeatmap.BeatmapsetId, &returnBeatmap.CreatorId, &returnBeatmap.Filename, &returnBeatmap.BeatmapMd5, &returnBeatmap.Version, &returnBeatmap.TotalLength, &returnBeatmap.DrainTime, &returnBeatmap.CountObjects, &returnBeatmap.CountNormal, &returnBeatmap.CountSlider, &returnBeatmap.CountSpinner, &returnBeatmap.DiffHp, &returnBeatmap.DiffCs, &returnBeatmap.DiffOd, &returnBeatmap.DiffStars, &returnBeatmap.Playmode, &returnBeatmap.RankingStatus, &returnBeatmap.LastUpdate, &returnBeatmap.SubmitDate, &returnBeatmap.ApproveDate, &returnBeatmap.BeatmapSource)

		beatmapQuery.Close()

		if scanErr != nil {
			return -2, Beatmap{}
		}

		return 0, returnBeatmap
	} else {
		beatmapQuery.Close()
		return -1, Beatmap{}
	}
}

func BeatmapsGetById(beatmapId int32) (queryResult int8, beatmap Beatmap) {
	beatmapQuery, beatmapQueryErr := Database.Query("SELECT beatmap_id, beatmapset_id, creator_id, filename, beatmap_md5, version, total_length, drain_time, count_objects, count_normal, count_slider, count_spinner, diff_hp, diff_cs, diff_od, diff_stars, playmode, ranking_status, last_update, submit_date, approve_date, beatmap_source FROM waffle.beatmaps WHERE beatmap_id = ?", beatmapId)

	if beatmapQueryErr != nil {
		if beatmapQuery != nil {
			beatmapQuery.Close()
		}

		return -2, Beatmap{}
	}

	if beatmapQuery.Next() {
		returnBeatmap := Beatmap{}

		scanErr := beatmapQuery.Scan(&returnBeatmap.BeatmapId, &returnBeatmap.BeatmapsetId, &returnBeatmap.CreatorId, &returnBeatmap.Filename, &returnBeatmap.BeatmapMd5, &returnBeatmap.Version, &returnBeatmap.TotalLength, &returnBeatmap.DrainTime, &returnBeatmap.CountObjects, &returnBeatmap.CountNormal, &returnBeatmap.CountSlider, &returnBeatmap.CountSpinner, &returnBeatmap.DiffHp, &returnBeatmap.DiffCs, &returnBeatmap.DiffOd, &returnBeatmap.DiffStars, &returnBeatmap.Playmode, &returnBeatmap.RankingStatus, &returnBeatmap.LastUpdate, &returnBeatmap.SubmitDate, &returnBeatmap.ApproveDate, &returnBeatmap.BeatmapSource)

		beatmapQuery.Close()

		if scanErr != nil {
			return -2, Beatmap{}
		}

		return 0, returnBeatmap
	} else {
		beatmapQuery.Close()
		return -1, Beatmap{}
	}
}
