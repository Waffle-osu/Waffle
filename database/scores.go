package database

import (
	"fmt"
	"strconv"
)

type Score struct {
	ScoreId         uint64
	BeatmapId       int
	BeatmapsetId    int
	UserId          uint64
	Playmode        int8
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

func (score Score) ScoresFormatLeaderboardScore(username string, onlineRank int32) string {
	perfectString := "0"

	if score.Perfect == 1 {
		perfectString = "1"
	}

	scoreIdstring := ""

	if onlineRank != 0 {
		scoreIdstring = strconv.FormatInt(int64(onlineRank), 10)
	}

	return fmt.Sprintf("%d|%s|%d|%d|%d|%d|%d|%d|%d|%d|%s|%d|%d|%s|%s\n", score.ScoreId, username, score.Score, score.MaxCombo, score.Hit50, score.Hit100, score.Hit300, score.HitMiss, score.HitKatu, score.HitKatu, perfectString, score.EnabledMods, score.UserId, scoreIdstring, score.Date)
}

func ScoresGetUserLeaderboardBest(beatmapId int32, userId uint64, mode int8) (queryResult int8, score Score, username string, onlineRank int64) {
	scoreQuery, scoreQueryErr := Database.Query("SELECT * FROM (SELECT ROW_NUMBER() OVER (ORDER BY score DESC) AS 'online_rank',  users.username, scores.* FROM waffle.scores LEFT JOIN waffle.users ON scores.user_id = users.user_id WHERE beatmap_id = ? AND leaderboard_best = 1 AND passed = 1 AND playmode = ? ORDER BY score DESC) t WHERE user_id = ?", beatmapId, mode, userId)

	if scoreQueryErr != nil {
		if scoreQuery != nil {
			scoreQuery.Close()
		}

		return -2, Score{}, "", -1
	}

	if scoreQuery.Next() {
		returnScore := Score{}

		var username string
		var onlineRank int64

		scanErr := scoreQuery.Scan(&onlineRank, &username, &returnScore.ScoreId, &returnScore.BeatmapId, &returnScore.BeatmapsetId, &returnScore.UserId, &returnScore.Playmode, &returnScore.Score, &returnScore.MaxCombo, &returnScore.Ranking, &returnScore.Hit300, &returnScore.Hit100, &returnScore.Hit50, &returnScore.HitMiss, &returnScore.HitGeki, &returnScore.HitKatu, &returnScore.EnabledMods, &returnScore.Perfect, &returnScore.Passed, &returnScore.Date, &returnScore.LeaderboardBest, &returnScore.MapsetBest, &returnScore.ScoreHash)

		scoreQuery.Close()

		if scanErr != nil {
			return -2, Score{}, "", -1
		}

		return 0, returnScore, username, onlineRank
	} else {
		scoreQuery.Close()
		return -1, Score{}, "", -1
	}
}

func ScoresGetBeatmapsetBestUserScore(beatmapsetId int32, userId uint64, mode int8) (queryResult int8, score Score) {
	scoreQuery, scoreQueryErr := Database.Query("SELECT score_id, beatmap_id, beatmapset_id, user_id, score, max_combo, ranking, hit300, hit100, hit50, hitMiss, hitGeki, hitKatu, enabled_mods, perfect, passed, date, leaderboard_best, mapset_best, score_hash, playmode FROM waffle.scores WHERE beatmapset_id = ? AND user_id = ? AND playmode = ? AND mapset_best = 1", beatmapsetId, userId, mode)

	if scoreQueryErr != nil {
		if scoreQuery != nil {
			scoreQuery.Close()
		}

		return -2, Score{}
	}

	if scoreQuery.Next() {
		returnScore := Score{}

		scanErr := scoreQuery.Scan(&returnScore.ScoreId, &returnScore.BeatmapId, &returnScore.BeatmapsetId, &returnScore.UserId, &returnScore.Score, &returnScore.MaxCombo, &returnScore.Ranking, &returnScore.Hit300, &returnScore.Hit100, &returnScore.Hit50, &returnScore.HitMiss, &returnScore.HitGeki, &returnScore.HitKatu, &returnScore.EnabledMods, &returnScore.Perfect, &returnScore.Passed, &returnScore.Date, &returnScore.LeaderboardBest, &returnScore.MapsetBest, &returnScore.ScoreHash, &returnScore.Playmode)

		scoreQuery.Close()

		if scanErr != nil {
			return -2, Score{}
		}

		return 0, returnScore
	} else {
		scoreQuery.Close()
		return -1, Score{}
	}
}

func ScoresGetBeatmapLeaderboardPlace(scoreId uint64, beatmapId int32) (queryResult int8, leaderboardPlace int64) {
	scoreQuery, scoreQueryErr := Database.Query("SELECT * FROM (SELECT score_id, beatmap_id, ROW_NUMBER() OVER (ORDER BY score DESC) AS 'rank' FROM waffle.scores WHERE beatmap_id = ? AND leaderboard_best = 1) t WHERE score_id = ?", beatmapId, scoreId)

	if scoreQueryErr != nil {
		if scoreQuery != nil {
			scoreQuery.Close()
		}

		return -2, -1
	}

	if scoreQuery.Next() {
		var scoreId uint64
		var beatmapId int32
		var leaderboardPlace int64

		scanErr := scoreQuery.Scan(&scoreId, &beatmapId, &leaderboardPlace)

		scoreQuery.Close()

		if scanErr != nil {
			return -2, -1
		}

		return 0, leaderboardPlace
	} else {
		scoreQuery.Close()
		return -1, -1
	}

}
