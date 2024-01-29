package actions

import (
	"Waffle/database"
	"errors"
	"fmt"
)

type GetLeaderboardsRequest struct {
	BeatmapChecksum string
	Filename        string
	Playmode        byte
	RankType        byte
	Mods            int32
	RequesterUserId int32

	GetScores        bool
	GetOffset        bool
	GetRating        bool
	GetRequesterBest bool
}

type LeaderboardScore struct {
	ScoreId    int64
	OnlineRank int64
	Username   string
	UserId     int64
	Score      int64

	Hit300  int32
	Hit100  int32
	Hit50   int32
	HitGeki int32
	HitKatu int32
	HitMiss int32

	MaxCombo int32
	Perfect  int8
	Mods     int32
	Date     string
}

func LeaderboardScoreFromDbScore(score database.Score, onlineRank int64, username string) LeaderboardScore {
	return LeaderboardScore{
		ScoreId:    int64(score.ScoreId),
		OnlineRank: onlineRank,
		Username:   username,
		UserId:     int64(score.UserId),
		Score:      int64(score.Score),
		Hit300:     int32(score.Hit300),
		Hit100:     int32(score.Hit100),
		Hit50:      int32(score.Hit50),
		HitGeki:    int32(score.HitGeki),
		HitKatu:    int32(score.HitKatu),
		HitMiss:    int32(score.HitMiss),
		MaxCombo:   int32(score.MaxCombo),
		Perfect:    score.Perfect,
		Mods:       int32(score.EnabledMods),
		Date:       score.Date,
	}
}

type GetLeaderboardsResponse struct {
	Error error

	HasOsz2Version   bool
	SubmissionStatus int8
	BeatmapId        int32
	BeatmapSetId     int32
	TotalScores      int32

	OnlineOffset int32
	DisplayTitle string
	OnlineRating float64

	PersonalBest LeaderboardScore
	Scores       []LeaderboardScore
}

func GetLeaderboards(request GetLeaderboardsRequest) GetLeaderboardsResponse {
	response := GetLeaderboardsResponse{
		Error:          nil,
		HasOsz2Version: false,
	}

	beatmap := database.Beatmap{}
	beatmapQueryResult := int8(0)

	beatmapset := database.Beatmapset{}
	beatmapsetQueryResult := int8(0)

	//Early 2007 clients don't send a filename.
	if request.Filename != "" {
		beatmapQueryResult, beatmap = database.BeatmapsGetByFilename(request.Filename)
	} else {
		beatmapQueryResult, beatmap = database.BeatmapsGetByMd5(request.BeatmapChecksum)
	}

	switch beatmapQueryResult {
	//Query failed
	case -2:
		response.Error = errors.New("beatmap query failed")

		fallthrough
	//Not found
	case -1:
		response.SubmissionStatus = -1
		return response
	}

	//Map needs updating
	if request.BeatmapChecksum != beatmap.BeatmapMd5 {
		response.SubmissionStatus = 1

		return response
	}

	beatmapsetQueryResult, beatmapset = database.BeatmapsetsGetBeatmapsetById(beatmap.BeatmapsetId)

	if beatmapsetQueryResult != 0 {
		response.Error = errors.New("beatmapset query failed")

		return response
	}

	response.BeatmapId = beatmap.BeatmapId
	response.BeatmapSetId = beatmap.BeatmapsetId
	response.DisplayTitle = fmt.Sprintf("[bold:0,size:20]%s|%s\n", beatmapset.Artist, beatmapset.Title)

	switch beatmap.RankingStatus {
	//Unranked
	case 0:
		response.SubmissionStatus = 0

		return response
	//Ranked
	case 1:
		response.SubmissionStatus = 2
	//Approved
	case 2:
		response.SubmissionStatus = 3
	}

	//Only a thing since late 2008
	if request.GetOffset {
		offsetQueryResult, offset := database.BeatmapOffsetsGetBeatmapOffset(beatmap.BeatmapId)

		if offsetQueryResult == 0 {
			response.OnlineOffset = offset.Offset
		}
	}

	//Only a thing since late 2008 aswell
	if request.GetRequesterBest && request.RequesterUserId != 0 {
		userBestScoreQueryResult, userBestScore, userUsername, userOnlineRank := database.ScoresGetUserLeaderboardBest(beatmap.BeatmapId, uint64(request.RequesterUserId), int8(request.Playmode))

		switch userBestScoreQueryResult {
		case -2:
			response.Error = errors.New("user best query failed")

			return response
		case -1:
			response.PersonalBest = LeaderboardScore{
				ScoreId: -1,
			}
		case 0:
			response.PersonalBest = LeaderboardScoreFromDbScore(userBestScore, userOnlineRank, userUsername)
		}
	}

	leaderboardQuery, leaderboardQueryErr := database.Database.Query("SELECT ROW_NUMBER() OVER (ORDER BY score DESC) AS 'online_rank', users.username, scores.* FROM waffle.scores LEFT JOIN waffle.users ON scores.user_id = users.user_id WHERE beatmap_id = ? AND leaderboard_best = 1 AND passed = 1 AND playmode = ? ORDER BY score DESC", beatmap.BeatmapId, int8(request.Playmode))

	if leaderboardQueryErr != nil {
		if leaderboardQuery != nil {
			leaderboardQuery.Close()
		}

		response.Error = leaderboardQueryErr

		return response
	}

	for leaderboardQuery.Next() {
		returnScore := database.Score{}

		var username string
		var onlineRank int64

		scanErr := leaderboardQuery.Scan(&onlineRank, &username, &returnScore.ScoreId, &returnScore.BeatmapId, &returnScore.BeatmapsetId, &returnScore.UserId, &returnScore.Playmode, &returnScore.Score, &returnScore.MaxCombo, &returnScore.Ranking, &returnScore.Hit300, &returnScore.Hit100, &returnScore.Hit50, &returnScore.HitMiss, &returnScore.HitGeki, &returnScore.HitKatu, &returnScore.EnabledMods, &returnScore.Perfect, &returnScore.Passed, &returnScore.Date, &returnScore.LeaderboardBest, &returnScore.MapsetBest, &returnScore.ScoreHash, &returnScore.Version)

		if scanErr != nil {
			leaderboardQuery.Close()

			response.Error = scanErr

			return response
		}

		response.Scores = append(response.Scores, LeaderboardScoreFromDbScore(returnScore, onlineRank, username))
	}

	leaderboardQuery.Close()

	return response
}
