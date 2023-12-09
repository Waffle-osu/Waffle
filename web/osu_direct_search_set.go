package web

import (
	"Waffle/bancho/client_manager"
	"Waffle/database"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func HandleDirectSearchSet(ctx *gin.Context) {
	beatmapSetId := ctx.Query("s")
	beatmapId := ctx.Query("b")
	threadId := ctx.Query("t")
	postId := ctx.Query("p")
	beatmapHash := ctx.Query("c")
	user := ctx.Query("u")

	if user == "" {
		ctx.String(http.StatusUnauthorized, "")

		return
	}

	//to make sure only one query ever gets triggered
	matchFound := false

	var beatmapSet database.Beatmapset
	var beatmap database.Beatmap

	if beatmapSetId != "" {
		setId, err := strconv.ParseInt(beatmapSetId, 10, 64)

		if err != nil {
			ctx.String(http.StatusBadRequest, "")

			return
		}

		queryResult, foundSet := database.BeatmapsetsGetBeatmapsetById(int32(setId))

		if queryResult != 0 {
			ctx.String(http.StatusOK, "")
		}

		beatmapSet = foundSet
		matchFound = true
	}

	if !matchFound && beatmapId != "" {
		mapId, err := strconv.ParseInt(beatmapId, 10, 64)

		if err != nil {
			ctx.String(http.StatusBadRequest, "")

			return
		}

		queryResultBeatmap, foundBeatmap := database.BeatmapsGetById(int32(mapId))

		if queryResultBeatmap != 0 {
			ctx.String(http.StatusOK, "")

			return
		}

		queryResultSet, foundSet := database.BeatmapsetsGetBeatmapsetById(foundBeatmap.BeatmapsetId)

		if queryResultSet != 0 {
			ctx.String(http.StatusOK, "")

			return
		}

		beatmapSet = foundSet
		beatmap = foundBeatmap
		matchFound = true
	}

	if !matchFound && beatmapHash != "" {
		queryResultBeatmap, foundBeatmap := database.BeatmapsGetByMd5(beatmapHash)

		if queryResultBeatmap != 0 {
			ctx.String(http.StatusOK, "")

			return
		}

		queryResultSet, foundSet := database.BeatmapsetsGetBeatmapsetById(foundBeatmap.BeatmapsetId)

		if queryResultSet != 0 {
			ctx.String(http.StatusOK, "")

			return
		}

		beatmapSet = foundSet
		beatmap = foundBeatmap
		matchFound = true
	}

	if !matchFound && threadId != "" {
		client := client_manager.GetClientByName(user)

		if client != nil {
			client.BanchoAnnounce("Due to Waffle not having a Forum, and not being able to conviniently query the osu! forums, osu!direct pickups via topic links don't work.")
		}

		ctx.String(http.StatusOK, "")

		return
	}

	if !matchFound && postId != "" {
		client := client_manager.GetClientByName(user)

		if client != nil {
			client.BanchoAnnounce("Due to Waffle not having a Forum, and not being able to conviniently query the osu! forums, osu!direct pickups via post links don't work.")
		}

		ctx.String(http.StatusOK, "")

		return
	}

	if beatmapSet.BeatmapsetId != 0 {
		returnFilename := fmt.Sprintf("%d %s - %s (%s).osz", beatmapSet.BeatmapsetId, beatmapSet.Artist, beatmapSet.Title, beatmapSet.Creator)

		//This is my way of determening the whole sets ranking status
		//ranked unless everythings approved
		sql := `
			SELECT * FROM 
				(SELECT SUM(ranking_status) FROM beatmaps WHERE beatmapset_id = ?) result UNION ALL
				(SELECT COUNT(*) FROM beatmaps WHERE beatmapset_id = ?) 
		`

		//row 1 = summed ranking status, if above or equal the diff count (row 2) means ranked, if equal to diff count * 2, approved
		//row 2 = diff count

		setStatusQuery, setStatusQueryErr := database.Database.Query(sql, beatmapSet.BeatmapsetId, beatmapSet.BeatmapsetId)

		if setStatusQueryErr != nil {
			ctx.String(http.StatusOK, "")

			return
		}

		var statusSum int64
		var diffCount int64

		//Status sum, has to be there so we do
		//a bunch of error checks that quit the
		//request if failed
		if setStatusQuery.Next() {
			scanErr := setStatusQuery.Scan(&statusSum)

			if scanErr != nil {
				ctx.String(http.StatusOK, "")

				return
			}
		} else {
			ctx.String(http.StatusOK, "")

			return
		}

		//Same here for the diff count
		if setStatusQuery.Next() {
			scanErr := setStatusQuery.Scan(&diffCount)

			if scanErr != nil {
				ctx.String(http.StatusOK, "")

				return
			}
		} else {
			ctx.String(http.StatusOK, "")

			return
		}

		//Math, since 1 means ranked and 2 means approved i can do this
		//say 4 diffs are approved, means all their statuses are 2, and summed up its 8
		//dividing by the diff count (4) returns 2
		//anything else returns 1.75 or smth in that range
		//works for pending too, since all their statuses are 0, and 0/4 is 0
		status := math.Floor(float64(statusSum) / float64(diffCount))

		rating := database.BeatmapRatingsGetBeatmapRating(beatmap.BeatmapsetId)

		fileSize := int64(0)

		if fileStats, err := os.Stat("oszs/" + strconv.FormatInt(int64(beatmap.BeatmapsetId), 10) + ".osz"); errors.Is(err, os.ErrNotExist) {

		} else {
			fileSize = fileStats.Size()
		}

		//Get last update
		lastUpdateQuery, lastUpdateQueryErr := database.Database.Query("SELECT beatmapset_id, last_update FROM beatmaps WHERE beatmapset_id = ? ORDER BY last_update DESC LIMIT 1", beatmapSet.BeatmapsetId)

		if lastUpdateQueryErr != nil {
			ctx.String(http.StatusOK, "")

			return
		}

		var lastUpdate time.Time

		if lastUpdateQuery.Next() {
			var queryTime string
			var queriedId int64

			lastUpdateQuery.Scan(&queriedId, &queryTime)

			date, dateErr := time.Parse("2006-01-02 15:04:05", queryTime)

			if dateErr != nil {
				ctx.String(http.StatusOK, "")

				return
			}

			lastUpdate = date
		} else {
			ctx.String(http.StatusOK, "")

			return
		}

		result := fmt.Sprintf("%s|%s|%s|%s|%d|%.2f|%s|%d|%d|%d|%d|%d|%d", returnFilename, beatmapSet.Artist, beatmapSet.Title, beatmapSet.Creator, int(status), rating, lastUpdate, beatmapSet.BeatmapsetId, 0, beatmapSet.HasVideo, beatmapSet.HasStoryboard, fileSize, fileSize)

		ctx.String(http.StatusOK, result)
	}
}
