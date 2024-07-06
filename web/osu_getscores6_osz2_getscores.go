package web

import (
	"Waffle/web/actions"
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
	osz2hash := ctx.Query("h")

	osz2client := false

	if beatmapsetId != "" || osz2hash != "" {
		osz2client = true
	}

	userId, parseErr := strconv.ParseInt(queryUserId, 10, 64)

	if parseErr != nil {
		ctx.String(http.StatusOK, "deliberatly fucked string to give out an error client side cuz pain\n")
		return
	}

	playmode, parseErr := strconv.ParseInt(queryPlaymode, 10, 64)

	if parseErr != nil {
		ctx.String(http.StatusOK, "deliberatly fucked string to give out an error client side cuz pain\n")
		return
	}

	leaderboardResponse := actions.GetLeaderboards(actions.GetLeaderboardsRequest{
		BeatmapChecksum:  beatmapChecksum,
		Filename:         osuFilename,
		Playmode:         byte(playmode),
		GetScores:        skipScores == "0",
		GetOffset:        true,
		GetRating:        true,
		GetRequesterBest: true,
		RequesterUserId:  int32(userId),
	})

	if leaderboardResponse.Error != nil {
		ctx.String(500, "string to make the client go aaaaaaaaaaaaaaaaaaaaaa")

		return
	}

	returnString := fmt.Sprintf("%d", leaderboardResponse.SubmissionStatus)

	if osz2client {
		returnString += fmt.Sprintf("|%t", leaderboardResponse.HasOsz2Version)
	}

	returnString += "\n"

	//Offset
	returnString += fmt.Sprintf("%d\n", leaderboardResponse.OnlineOffset)

	//Display Title
	returnString += leaderboardResponse.DisplayTitle

	//Online Rating
	returnString += fmt.Sprintf("%.2f\n", leaderboardResponse.OnlineRating)

	if skipScores == "1" {
		ctx.String(http.StatusOK, returnString)
		return
	}

	formatScore := func(score actions.LeaderboardScore) string {
		perfectString := "0"

		if score.Perfect == 1 {
			perfectString = "1"
		}

		scoreIdstring := ""

		if score.OnlineRank != 0 {
			scoreIdstring = strconv.FormatInt(score.OnlineRank, 10)
		}

		return fmt.Sprintf("%d|%s|%d|%d|%d|%d|%d|%d|%d|%d|%s|%d|%d|%s|%s\n", score.ScoreId, score.Username, score.Score, score.MaxCombo, score.Hit50, score.Hit100, score.Hit300, score.HitMiss, score.HitKatu, score.HitGeki, perfectString, score.Mods, score.UserId, scoreIdstring, score.Date)
	}

	//0 is invalid and -1 is too
	if leaderboardResponse.PersonalBest.ScoreId > 0 {
		returnString += formatScore(leaderboardResponse.PersonalBest)
	} else {
		returnString += "\n"
	}

	for _, score := range leaderboardResponse.Scores {
		returnString += formatScore(score)
	}

	ctx.String(http.StatusOK, returnString)
}
