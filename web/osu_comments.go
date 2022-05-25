package web

import (
	"Waffle/bancho/packets"
	"Waffle/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	BeatmapCommentsTargetNone   = 0
	BeatmapCommentsTargetMap    = 1
	BeatmapCommentsTargetSong   = 2
	BeatmapCommentsTargetReplay = 3
)

func HandleOsuComments(ctx *gin.Context) {
	queryUsername := ctx.PostForm("u")
	queryPassword := ctx.PostForm("p")
	queryBeatmapId := ctx.PostForm("b")
	queryBeatmapSetId := ctx.PostForm("s")
	queryScoreId := ctx.PostForm("r")
	queryAction := ctx.PostForm("a")

	//post only
	queryTarget := ctx.PostForm("target")
	queryStartTime := ctx.PostForm("startttime")
	queryComment := ctx.PostForm("comment")

	beatmapId, parseErr1 := strconv.ParseInt(queryBeatmapId, 10, 64)
	beatmapsetId, parseErr2 := strconv.ParseInt(queryBeatmapSetId, 10, 64)
	scoreId, parseErr3 := strconv.ParseInt(queryScoreId, 10, 64)

	//post only
	startTime := int32(0)
	target := int8(0)

	if queryStartTime != "" {
		startTime, parseErr := strconv.ParseInt(queryStartTime, 10, 64)

		if parseErr != nil {
			ctx.String(http.StatusBadRequest, "")
			return
		}
	}

	if queryTarget != "" {
		switch queryTarget {
		case "none":
			target = 0
		case "map":
			target = 1
		case "song":
			target = 2
		case "replay":
			target = 3
		}
	}

	if parseErr1 != nil || parseErr2 != nil || parseErr3 != nil {
		ctx.String(http.StatusBadRequest, "")
		return
	}

	userId, authResult := database.AuthenticateUser(queryUsername, queryPassword)

	if userId == -2 {
		ctx.String(http.StatusInternalServerError, "")
		return
	}

	if !authResult {
		ctx.String(http.StatusUnauthorized, "")
		return
	}

	_, databaseUser := database.UserFromDatabaseById(uint64(userId))

	switch queryAction {
	case "get":
	case "post":
		formatString := ""

		playerQuery, playerQueryErr := database.Database.Query("SELECT score_id, user_id FROM waffle.scores WHERE score_id = ?", uint64(scoreId))

		if playerQueryErr != nil {
			if playerQuery != nil {
				playerQuery.Close()
			}

			ctx.String(http.StatusInternalServerError, "")
			return
		}

		if playerQuery.Next() {
			var foundScoreId, foundUserId uint64

			scanErr := playerQuery.Scan(&foundScoreId, &foundUserId)

			if scanErr != nil {
				ctx.String(http.StatusInternalServerError, "")
				return
			}

			if foundUserId == uint64(userId) {
				formatString = "player"
			}

			playerQuery.Close()
		}

		creatorQuery, creatorQueryErr := database.Database.Query("SELECT beatmapset_id, creator_id, creator FROM waffle.beatmapsets WHERE beatmapset_id = ?", int32(beatmapId))

		if creatorQueryErr != nil {
			if creatorQuery != nil {
				creatorQuery.Close()
			}

			ctx.String(http.StatusInternalServerError, "")
			return
		}

		if creatorQuery.Next() {
			var foundSetId int32
			var foundCreatorId int64
			var foundCreator string

			scanErr := creatorQuery.Scan(&foundSetId, &foundCreatorId, &foundCreator)

			if scanErr != nil {
				ctx.String(http.StatusInternalServerError, "")
				return
			}

			if foundCreatorId == int64(userId) || databaseUser.Username == foundCreator {
				formatString = "creator"
			}

			creatorQuery.Close()
		}

		if (databaseUser.Privileges & packets.UserPermissionsBAT) > 0 {
			formatString = "bat"
		}

		if (databaseUser.Privileges & packets.UserPermissionsSupporter) > 0 {
			formatString = "subscriber"
		}

		insertQuery, insertQueryErr := database.Database.Query("INSERT INTO waffle.beatmap_comments (user_id, beatmap_id, beatmapset_id, score_id, time, target, comment, format_string) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", uint64(userId), int32(beatmapId), int32(beatmapsetId), uint64(scoreId), startTime, target, queryComment, formatString)

		if insertQueryErr != nil {
			if insertQuery != nil {
				insertQuery.Close()
			}

			ctx.String(http.StatusInternalServerError, "")
			return
		}

		insertQuery.Close()
	default:
		ctx.String(http.StatusBadRequest, "")
		return
	}
}
