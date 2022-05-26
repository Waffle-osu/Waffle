package web

import (
	"Waffle/database"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HandleOsuGetLeaderboards(ctx *gin.Context) {
	skipScores := ctx.Query("s")
	beatmapChecksum := ctx.Query("c")
	osuFilename := ctx.Query("f")
	queryUserId := ctx.Query("u")
	queryPlaymode := ctx.Query("m")
	beatmapsetId := ctx.Query("i")
	//osz2hash := ctx.Query("h")

	osz2client := false

	if beatmapsetId != "" {
		osz2client = true
	}

	userId, parseErr := strconv.ParseInt(queryUserId, 10, 64)
	//playmode, parseErr := strconv.ParseInt(queryPlaymode, 10, 64)

	if parseErr != nil {
		ctx.String(http.StatusOK, "deliberatly fucked string to give out an error client side cuz pain\n")
		return
	}

	leaderboardBeatmapQueryResult, leaderboardBeatmap := database.BeatmapsGetByFilename(osuFilename)

	if leaderboardBeatmapQueryResult == -2 {
		ctx.String(http.StatusOK, "deliberatly fucked string to give out an error client side cuz pain\n")
		return
	}

	if leaderboardBeatmapQueryResult == -1 {
		response := "-1"

		if osz2client {
			response += "|false"
		}

		ctx.String(http.StatusOK, response)
		return
	}

	if beatmapChecksum != leaderboardBeatmap.BeatmapMd5 {
		response := "1"

		if osz2client {
			response += "|false"
		}

		ctx.String(http.StatusOK, response)
		return
	}

	beatmapsetQueryResult, beatmapset := database.BeatmapsetsGetBeatmapsetById(leaderboardBeatmap.BeatmapsetId)

	if beatmapsetQueryResult == -2 {
		ctx.String(http.StatusOK, "deliberatly fucked string to give out an error client side cuz pain\n")
		return
	}

	returnString := ""

	returnRankedStatus := "0"

	switch leaderboardBeatmap.RankingStatus {
	case 0:
		if osz2client {
			ctx.String(http.StatusOK, "0|false")
		} else {
			ctx.String(http.StatusOK, "0")
		}
		return
	case 1:
		returnRankedStatus = "2"
	case 2:
		returnRankedStatus = "3"
	}
	//Ranked Status|Server has osz2 of map
	returnString += returnRankedStatus

	if osz2client {
		returnString += "|false\n"
	} else {
		returnString += "\n"
	}

	//Online Offset, currently we don't store any so eh, TODO
	returnString += "0\n"
	//Display Title
	returnString += fmt.Sprintf("[bold:0,size:20]%s|%s\n", beatmapset.Artist, beatmapset.Title)

	getRatingInfoQuery, getRatingInfoQueryErr := database.Database.Query("SELECT * FROM beatmap_ratings WHERE beatmapset_id = ?", beatmapset.BeatmapsetId)

	if getRatingInfoQueryErr != nil {
		if getRatingInfoQuery != nil {
			getRatingInfoQuery.Close()
		}

		ctx.String(http.StatusOK, "because server fucked up")
		return
	}

	var ratingSum, votes int64

	if getRatingInfoQuery.Next() {
		var beatmapsetId int32

		scanErr := getRatingInfoQuery.Scan(&beatmapsetId, &ratingSum, &votes)

		if getRatingInfoQuery != nil {
			getRatingInfoQuery.Close()
		}

		if scanErr != nil {
			ctx.String(http.StatusOK, "because server fucked up")
			return
		}
	}

	if votes == 0 {
		votes++
	}

	totalRating := float64(ratingSum) / float64(votes)

	//Online Rating, currently rating doesnt exist, so TODO
	returnString += fmt.Sprintf("%.2f\n", totalRating)

	if skipScores == "1" {
		ctx.String(http.StatusOK, returnString)
		return
	}

	playmode, parseErr := strconv.ParseInt(queryPlaymode, 10, 64)

	if parseErr != nil {
		ctx.String(http.StatusOK, "deliberatly fucked string to give out an error client side cuz pain\n")
		return
	}

	userBestScoreQueryResult, userBestScore, userUsername, userOnlineRank := database.ScoresGetUserLeaderboardBest(leaderboardBeatmap.BeatmapId, uint64(userId), int8(playmode))

	if userBestScoreQueryResult == -1 || userBestScore.Passed == 0 {
		returnString += "\n"
	} else {
		returnString += userBestScore.ScoresFormatLeaderboardScore(userUsername, int32(userOnlineRank))
	}

	leaderboardQuery, leaderboardQueryErr := database.Database.Query("SELECT ROW_NUMBER() OVER (ORDER BY score DESC) AS 'online_rank', users.username, scores.* FROM waffle.scores LEFT JOIN waffle.users ON scores.user_id = users.user_id WHERE beatmap_id = ? AND leaderboard_best = 1 AND passed = 1 AND playmode = ? ORDER BY score DESC", leaderboardBeatmap.BeatmapId, int8(playmode))

	if leaderboardQueryErr != nil {
		if leaderboardQuery != nil {
			leaderboardQuery.Close()
		}

		ctx.String(http.StatusOK, "deliberatly fucked string to give out an error client side cuz pain\n")
		return
	}

	for leaderboardQuery.Next() {
		returnScore := database.Score{}

		var username string
		var onlineRank int64

		scanErr := leaderboardQuery.Scan(&onlineRank, &username, &returnScore.ScoreId, &returnScore.BeatmapId, &returnScore.BeatmapsetId, &returnScore.UserId, &returnScore.Playmode, &returnScore.Score, &returnScore.MaxCombo, &returnScore.Ranking, &returnScore.Hit300, &returnScore.Hit100, &returnScore.Hit50, &returnScore.HitMiss, &returnScore.HitGeki, &returnScore.HitKatu, &returnScore.EnabledMods, &returnScore.Perfect, &returnScore.Passed, &returnScore.Date, &returnScore.LeaderboardBest, &returnScore.MapsetBest, &returnScore.ScoreHash)

		if scanErr != nil {
			leaderboardQuery.Close()
			ctx.String(http.StatusOK, "deliberatly fucked string to give out an error client side cuz pain\n")
			return
		}

		returnString += returnScore.ScoresFormatLeaderboardScore(username, int32(onlineRank))
	}

	leaderboardQuery.Close()

	ctx.String(http.StatusOK, returnString)
}
