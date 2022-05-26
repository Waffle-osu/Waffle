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

func BeatmapsetsGetBeatmapsetById(beatmapsetId int32) (queryResult int8, beatmapset Beatmapset) {
	beatmapsetQuery, beatmapsetQueryErr := Database.Query("SELECT beatmapset_id, creator_id, artist, title, creator, source, tags, has_video, has_storyboard, bpm FROM beatmapsets WHERE beatmapset_id = ?", beatmapsetId)

	if beatmapsetQueryErr != nil {
		if beatmapsetQuery != nil {
			beatmapsetQuery.Close()
		}

		return -2, Beatmapset{}
	}

	if beatmapsetQuery.Next() {
		returnBeatmapset := Beatmapset{}

		scanErr := beatmapsetQuery.Scan(&returnBeatmapset.BeatmapsetId, &returnBeatmapset.CreatorId, &returnBeatmapset.Artist, &returnBeatmapset.Title, &returnBeatmapset.Creator, &returnBeatmapset.Source, &returnBeatmapset.Tags, &returnBeatmapset.HasVideo, &returnBeatmapset.HasStoryboard, &returnBeatmapset.Bpm)

		beatmapsetQuery.Close()

		if scanErr != nil {
			return -2, Beatmapset{}
		}

		return 0, returnBeatmapset
	} else {
		beatmapsetQuery.Close()
		return -1, Beatmapset{}
	}
}
