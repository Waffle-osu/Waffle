package web

import (
	"Waffle/bancho/packets"
	"Waffle/database"
	"fmt"
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
	queryStartTime := ctx.PostForm("starttime")
	queryComment := ctx.PostForm("comment")

	beatmapId, parseErr1 := strconv.ParseInt(queryBeatmapId, 10, 64)
	beatmapsetId, parseErr2 := strconv.ParseInt(queryBeatmapSetId, 10, 64)
	scoreId, parseErr3 := strconv.ParseInt(queryScoreId, 10, 64)

	//post only
	startTime := int32(0)
	target := int8(0)

	if queryStartTime != "" {
		startTimeParsed, parseErr := strconv.ParseInt(queryStartTime, 10, 64)

		if parseErr != nil {
			ctx.String(http.StatusBadRequest, "")
			return
		}

		startTime = int32(startTimeParsed)
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
		getSqlQuery := `
SELECT * FROM
	(SELECT * FROM waffle.beatmap_comments WHERE beatmap_id = ? AND target = 1) beatmapresults UNION ALL
	(SELECT * FROM waffle.beatmap_comments WHERE beatmapset_id = ? AND target = 2) UNION ALL
	(SELECT * FROM waffle.beatmap_comments WHERE score_id = ? AND target = 3)  
ORDER BY time ASC
		`

		getQuery, getQueryErr := database.Database.Query(getSqlQuery, int32(beatmapId), int32(beatmapsetId), uint64(scoreId))

		if getQueryErr != nil {
			if getQuery != nil {
				getQuery.Close()
			}

			ctx.String(http.StatusInternalServerError, "")
			return
		}

		returnString := ""

		for getQuery.Next() {
			comment := database.BeatmapComment{}

			scanErr := getQuery.Scan(&comment.CommentId, &comment.UserId, &comment.BeatmapId, &comment.BeatmapSetId, &comment.ScoreId, &comment.Time, &comment.Target, &comment.Comment, &comment.FormatString)

			if scanErr != nil {
				ctx.String(http.StatusInternalServerError, "")
				return
			}

			parsedTarget := "none"

			switch comment.Target {
			case BeatmapCommentsTargetNone:
				parsedTarget = "none"
			case BeatmapCommentsTargetReplay:
				parsedTarget = "replay"
			case BeatmapCommentsTargetSong:
				parsedTarget = "song"
			case BeatmapCommentsTargetMap:
				parsedTarget = "map"
			}

			//this is split by '\t'
			//[0]: time
			//[1]: Enum to String (None, Replay, Map, Song) but lowercase
			//[2]: Format String
			//[3]: Comment
			returnString += fmt.Sprintf("%s\t%s\t%s\t%s\n", strconv.FormatInt(int64(comment.Time), 10), parsedTarget, comment.FormatString, comment.Comment)
		}

		if getQuery != nil {
			getQuery.Close()
		}

		ctx.String(http.StatusOK, returnString)
		return
	case "post":
		formatString := ""

		//Check if it's the player sending the comment
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

		//check if the creator is sending the comment
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

		//check if its a osu!supporter sending the comment
		if (databaseUser.Privileges & packets.UserPermissionsSupporter) > 0 {
			formatString = "subscriber"
		}

		//check if its a BAT sending the comment
		if (databaseUser.Privileges & packets.UserPermissionsBAT) > 0 {
			formatString = "bat"
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

		ctx.String(http.StatusOK, "")
	default:
		ctx.String(http.StatusBadRequest, "")
		return
	}
}
