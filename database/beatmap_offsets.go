package database

type BeatmapOffset struct {
	OffsetId  int64
	BeatmapId int32
	Offset    int32
}

// Gets a beatmap offset for a map
func BeatmapOffsetsGetBeatmapOffset(beatmapId int32) (result int8, offset BeatmapOffset) {
	offsetQuery, offsetQueryErr := Database.Query("SELECT * FROM waffle.beatmap_offsets WHERE beatmap_id = ?", beatmapId)

	if offsetQueryErr != nil {
		if offsetQuery != nil {
			offsetQuery.Close()
		}

		return -2, BeatmapOffset{}
	}

	if offsetQuery.Next() {
		offset := BeatmapOffset{}

		scanErr := offsetQuery.Scan(&offset.OffsetId, &offset.BeatmapId, &offset.Offset)

		offsetQuery.Close()

		if scanErr != nil {
			return -2, BeatmapOffset{}
		}

		return 0, offset
	} else {
		return -1, BeatmapOffset{}
	}
}
