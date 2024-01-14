package bss

import (
	"Waffle/database"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HandlePost3(ctx *gin.Context) {
	username := ctx.PostForm("u")
	password := ctx.PostForm("p")
	setId := ctx.PostForm("b")
	subject := ctx.PostForm("subject")
	message := ctx.PostForm("message")
	complete := ctx.PostForm("complete")
	notify := ctx.PostForm("notify")

	userId, success := database.AuthenticateUser(username, password)

	if !success {
		ctx.String(401, "Authentication failed!")

		return
	}

	parsedSetId, parseErr := strconv.ParseInt(setId, 10, 64)

	if parseErr != nil {
		ctx.String(400, "Invalid Set ID")

		return
	}

	queryResult, beatmapSet := database.BeatmapsetsGetBeatmapsetById(int32(parsedSetId))

	if queryResult != 0 {
		ctx.String(500, "Internal queries failed")

		return
	}

	if (-beatmapSet.CreatorId) != int64(userId) {
		ctx.String(401, "Not your beatmap!")

		return
	}

	existsQuerySql := "SELECT COUNT(*) FROM osu_beatmap_posts WHERE beatmapset_id = ?"
	existsQuery, existsQueryErr := database.Database.Query(existsQuerySql, parsedSetId)

	if existsQueryErr != nil {
		ctx.String(500, "Internal queries failed")

		return
	}

	existsQueryResult := int64(0)

	existsQuery.Next()
	scanErr := existsQuery.Scan(&existsQueryResult)
	existsQuery.Close()

	if scanErr != nil {
		ctx.String(500, "Internal queries failed")

		return
	}

	actualNotify := byte(0)
	actualComplete := byte(0)

	if notify == "1" {
		actualNotify = 1
	}

	if complete == "1" {
		actualComplete = 1
	}

	if existsQueryResult > 0 {
		//Update
		updateSql := "UPDATE osu_beatmap_posts SET beatmapset_id = ?, subject = ?, message = ?, notify = ?, complete = ?"
		_, updateQueryErr := database.Database.Exec(updateSql, parsedSetId, subject, message, actualNotify, actualComplete)

		if updateQueryErr != nil {
			ctx.String(500, "Internal queries failed")

			return
		}
	} else {
		//Create
		createSql := "INSERT INTO osu_beatmap_posts (beatmapset_id, subject, message, notify, complete) VALUES (?, ?, ?, ?, ?)"
		_, createQueryErr := database.Database.Exec(createSql, parsedSetId, subject, message, actualNotify, actualComplete)

		if createQueryErr != nil {
			ctx.String(500, "Internal queries failed")

			return
		}

		ctx.String(200, fmt.Sprintf("%d", parsedSetId))
	}
}
