package web

import (
	"Waffle/database"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func HandleOsuOsz2GetScores(ctx *gin.Context) {
	//gettingScores := ctx.Query("s")
	beatmapChecksum := ctx.Query("c")
	osuFilename := ctx.Query("f")
	queryUserId := ctx.Query("u")
	//queryPlaymode := ctx.Query("m")
	//beatmapsetId := ctx.Query("i")
	//osz2hash := ctx.Query("h")

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
		ctx.String(http.StatusOK, "-1|false")
		return
	}

	if beatmapChecksum != leaderboardBeatmap.BeatmapMd5 {
		ctx.String(http.StatusOK, "1|false")
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
		ctx.String(http.StatusOK, "0|false")
		return
	case 1:
		returnRankedStatus = "2"
		break
	case 2:
		returnRankedStatus = "3"
		break
	}

	returnString += returnRankedStatus + "|false\n" //Ranked Status|Server has osz2 of map
	returnString += "0\n"                           //Online Offset, currently we don't store any so eh, TODO
	returnString += "0\n"                           //Online Rating, currently rating doesnt exist, so TODO

	returnString += fmt.Sprintf("[bold:0,size:20]%s|%s", beatmapset.Artist, beatmapset.Title) //Display Title

	userBestScoreQueryResult, userBestScore := database.ScoresGetUserLeaderboardBest(leaderboardBeatmap.BeatmapId, uint64(userId))

	if userBestScoreQueryResult == -1 || userBestScore.Passed == 0 {
		returnString += "\n"
	}

	//Currently, has to still be done,
	//but the first line has to have the Ranking Status seperated with a | and whether the server has a .osz2 file of the map
	ctx.String(http.StatusOK, "2|false")
}
