package database

type Score struct {
	ScoreId         uint64
	BeatmapId       int
	BeatmapsetId    int
	UserId          uint64
	Score           int
	MaxCombo        int
	Ranking         string
	Hit300          int
	Hit100          int
	Hit50           int
	HitMiss         int
	HitGeki         int
	HitKatu         int
	EnabledMods     int
	Perfect         int8
	Passed          int8
	Date            string
	LeaderboardBest int8
	MapsetBest      int8
	ScoreHash       string
}

func ScoresGetUserLeaderboardBest(beatmapId int32, userId uint64) (queryResult int8, score Score) {
	scoreQuery, scoreQueryErr := Database.Query("SELECT score_id, beatmap_id, beatmapset_id, user_id, score, max_combo, ranking, hit300, hit100, hit50, hitMiss, hitGeki, hitKatu, enabled_mods, perfect, passed, date, leaderboard_best, mapset_best, score_hash FROM waffle.scores WHERE beatmap_id = ? AND user_id = ? AND leaderboard_best = 1", beatmapId, userId)

	if scoreQueryErr != nil {
		if scoreQuery != nil {
			scoreQuery.Close()
		}

		return -2, Score{}
	}

	if scoreQuery.Next() {
		returnScore := Score{}

		scanErr := scoreQuery.Scan(&returnScore.ScoreId, &returnScore.BeatmapId, &returnScore.BeatmapsetId, &returnScore.UserId, &returnScore.Score, &returnScore.MaxCombo, &returnScore.Ranking, &returnScore.Hit300, &returnScore.Hit100, &returnScore.Hit50, &returnScore.HitMiss, &returnScore.HitGeki, &returnScore.HitKatu, &returnScore.EnabledMods, &returnScore.Perfect, &returnScore.Passed, &returnScore.Date, &returnScore.LeaderboardBest, &returnScore.MapsetBest, &returnScore.ScoreHash)

		if scanErr != nil {
			scoreQuery.Close()
			return -2, Score{}
		}

		scoreQuery.Close()
		return 0, returnScore
	} else {
		scoreQuery.Close()
		return -1, Score{}
	}
}

func ScoresGetBeatmapsetUserScore(beatmapsetId int32, userId uint64) (queryResult int8, score Score) {
	scoreQuery, scoreQueryErr := Database.Query("SELECT score_id, beatmap_id, beatmapset_id, user_id, score, max_combo, ranking, hit300, hit100, hit50, hitMiss, hitGeki, hitKatu, enabled_mods, perfect, passed, date, leaderboard_best, mapset_best, score_hash FROM waffle.scores WHERE beatmapset_id = ? AND user_id = ? AND mapset_best = 1", beatmapsetId, userId)

	if scoreQueryErr != nil {
		if scoreQuery != nil {
			scoreQuery.Close()
		}

		return -2, Score{}
	}

	if scoreQuery.Next() {
		returnScore := Score{}

		scanErr := scoreQuery.Scan(&returnScore.ScoreId, &returnScore.BeatmapId, &returnScore.BeatmapsetId, &returnScore.UserId, &returnScore.Score, &returnScore.MaxCombo, &returnScore.Ranking, &returnScore.Hit300, &returnScore.Hit100, &returnScore.Hit50, &returnScore.HitMiss, &returnScore.HitGeki, &returnScore.HitKatu, &returnScore.EnabledMods, &returnScore.Perfect, &returnScore.Passed, &returnScore.Date, &returnScore.LeaderboardBest, &returnScore.MapsetBest, &returnScore.ScoreHash)

		if scanErr != nil {
			scoreQuery.Close()
			return -2, Score{}
		}

		scoreQuery.Close()
		return 0, returnScore
	} else {
		scoreQuery.Close()
		return -1, Score{}
	}
}

func ScoresGetBeatmapLeaderboardPlace(scoreId uint64, beatmapId int32) (queryResult int8, leaderboardPlace int64) {
	scoreQuery, scoreQueryErr := Database.Query("SELECT score_id, beatmap_id, ROW_NUMBER() OVER (ORDER BY score DESC) AS 'rank' FROM waffle.scores WHERE beatmap_id = ? AND score_id = ? AND leaderboard_best = 1", beatmapId, scoreId)

	if scoreQueryErr != nil {
		if scoreQuery != nil {
			scoreQuery.Close()
		}

		return -2, -1
	}

	if scoreQuery.Next() {
		var leaderboardPlace int64

		scanErr := scoreQuery.Scan(&leaderboardPlace)

		if scanErr != nil {
			scoreQuery.Close()
			return -2, -1
		}

		scoreQuery.Close()
		return 0, leaderboardPlace
	} else {
		scoreQuery.Close()
		return -1, -1
	}

}
