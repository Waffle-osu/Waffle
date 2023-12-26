package database

// Gets the beatmap rating of a set
func BeatmapRatingsGetBeatmapRating(beatmapsetId int32) float64 {
	getRatingInfoQuery, getRatingInfoQueryErr := Database.Query("SELECT * FROM beatmap_ratings WHERE beatmapset_id = ?", beatmapsetId)

	if getRatingInfoQueryErr != nil {
		if getRatingInfoQuery != nil {
			getRatingInfoQuery.Close()
		}

		return 0
	}

	var ratingSum, votes int64

	if getRatingInfoQuery.Next() {
		var beatmapsetId int32

		scanErr := getRatingInfoQuery.Scan(&beatmapsetId, &ratingSum, &votes)

		getRatingInfoQuery.Close()

		if scanErr != nil {
			return 0
		}
	}

	if votes == 0 {
		votes++
	}

	return float64(ratingSum) / float64(votes)
}
