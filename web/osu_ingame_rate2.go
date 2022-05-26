package web

import (
	"Waffle/database"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HandleOsuIngameRate2(ctx *gin.Context) {
	queryUsername := ctx.Query("u")
	queryPassword := ctx.Query("p")
	queryBeatmapMd5 := ctx.Query("c")
	queryUserRating := ctx.Query("v")

	isSubmission := queryUserRating != ""

	userId, authResult := database.AuthenticateUser(queryUsername, queryPassword)

	if userId == -2 {
		ctx.String(http.StatusOK, "because server fucked up")
		return
	}

	if !authResult {
		ctx.String(http.StatusOK, "auth fail")
		return
	}

	_, userData := database.UserFromDatabaseById(uint64(userId))

	beatmapQueryResult, beatmap := database.BeatmapsGetByMd5(queryBeatmapMd5)

	if beatmapQueryResult == -2 {
		ctx.String(http.StatusOK, "because server fucked up")
		return
	}

	if beatmapQueryResult == -1 {
		ctx.String(http.StatusOK, "no exist")
		return
	}

	if beatmap.RankingStatus != 2 && beatmap.RankingStatus != 1 {
		ctx.String(http.StatusOK, "not ranked")
		return
	}

	beatmapsetQueryResult, beatmapset := database.BeatmapsetsGetBeatmapsetById(beatmap.BeatmapsetId)

	if beatmapsetQueryResult == -2 {
		ctx.String(http.StatusOK, "because server fucked up")
		return
	}

	if beatmapsetQueryResult == -1 {
		ctx.String(http.StatusOK, "no exist")
		return
	}

	getMapRatingsExistQuery, getMapRatingsExistQueryErr := database.Database.Query("SELECT COUNT(*) AS 'count' FROM waffle.beatmap_ratings WHERE beatmapset_id = ?", beatmapset.BeatmapsetId)

	if getMapRatingsExistQueryErr != nil {
		if getMapRatingsExistQuery != nil {
			getMapRatingsExistQuery.Close()
		}

		ctx.String(http.StatusOK, "because server fucked up")
		return
	}

	if getMapRatingsExistQuery.Next() {
		var count int64

		scanErr := getMapRatingsExistQuery.Scan(&count)

		getMapRatingsExistQuery.Close()

		if scanErr != nil {
			ctx.String(http.StatusOK, "because server fucked up")
			return
		}

		if count == 0 {
			insertBeatmapRatingTrackQuery, insertBeatmapRatingTrackQueryErr := database.Database.Query("INSERT INTO beatmap_ratings (beatmapset_id) VALUES (?)", beatmapset.BeatmapsetId)

			if insertBeatmapRatingTrackQueryErr != nil {
				if insertBeatmapRatingTrackQuery != nil {
					insertBeatmapRatingTrackQuery.Close()
				}

				ctx.String(http.StatusOK, "because server fucked up")
				return
			}

			if insertBeatmapRatingTrackQuery != nil {
				insertBeatmapRatingTrackQuery.Close()
			}
		}
	}

	if isSubmission {
		submittedRating, parseErr := strconv.ParseInt(queryUserRating, 10, 64)

		if parseErr != nil {
			ctx.String(http.StatusBadRequest, "")
			return
		}

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

			getRatingInfoQuery.Close()

			if scanErr != nil {
				ctx.String(http.StatusOK, "because server fucked up")
				return
			}
		}

		userHasSubmittedRatingQuery, userHasSubmittedRatingQueryErr := database.Database.Query("SELECT COUNT(*) AS 'count' FROM waffle.beatmap_ratings_submissions WHERE user_id = ? AND beatmapset_id = ?", uint64(userId), beatmapset.BeatmapsetId)

		if userHasSubmittedRatingQueryErr != nil {
			if userHasSubmittedRatingQuery != nil {
				userHasSubmittedRatingQuery.Close()
			}

			ctx.String(http.StatusOK, "because server fucked up")
			return
		}

		if userHasSubmittedRatingQuery.Next() {
			var count int64

			scanErr := userHasSubmittedRatingQuery.Scan(&count)

			userHasSubmittedRatingQuery.Close()

			if scanErr != nil {
				ctx.String(http.StatusOK, "because server fucked up")
				return
			}

			if count != 0 {
				ctx.String(http.StatusOK, fmt.Sprintf("%.2f", float64(ratingSum)/float64(votes)))
				return
			}
		}

		ratingSum += submittedRating
		votes++

		newRating := float64(ratingSum) / float64(votes)

		updateRatingInfoQuery, updateRatingInfoQueryErr := database.Database.Query("UPDATE waffle.beatmap_ratings SET rating_sum = ?, votes = ? WHERE beatmapset_id = ?", ratingSum, votes, beatmapset.BeatmapsetId)

		if updateRatingInfoQueryErr != nil {
			if updateRatingInfoQuery != nil {
				updateRatingInfoQuery.Close()
			}

			ctx.String(http.StatusInternalServerError, "because server fucked up")
			return
		}

		if updateRatingInfoQuery != nil {
			updateRatingInfoQuery.Close()
		}

		insertUserSubmittedQuery, insertUserSubmittedQueryErr := database.Database.Query("INSERT INTO waffle.beatmap_ratings_submissions (user_id, beatmapset_id) VALUES (?, ?)", uint64(userId), beatmapset.BeatmapsetId)

		if insertUserSubmittedQueryErr != nil {
			if insertUserSubmittedQuery != nil {
				insertUserSubmittedQuery.Close()
			}

			ctx.String(http.StatusInternalServerError, "because server fucked up")
			return
		}

		if insertUserSubmittedQuery != nil {
			insertUserSubmittedQuery.Close()
		}

		ctx.String(http.StatusOK, fmt.Sprintf("%.2f", newRating))
		return
	} else {
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

			getRatingInfoQuery.Close()

			if scanErr != nil {
				ctx.String(http.StatusOK, "because server fucked up")
				return
			}
		}

		if votes == 0 {
			votes++
		}

		totalRating := float64(ratingSum) / float64(votes)

		if beatmapset.CreatorId == int64(userId) || userData.Username == beatmapset.Creator {
			ctx.String(http.StatusOK, fmt.Sprintf("creator\n%.2f", totalRating))
			return
		}

		userHasSubmittedRatingQuery, userHasSubmittedRatingQueryErr := database.Database.Query("SELECT COUNT(*) AS 'count' FROM waffle.beatmap_ratings_submissions WHERE user_id = ? AND beatmapset_id = ?", uint64(userId), beatmapset.BeatmapsetId)

		if userHasSubmittedRatingQueryErr != nil {
			if userHasSubmittedRatingQuery != nil {
				userHasSubmittedRatingQuery.Close()
			}

			ctx.String(http.StatusOK, "because server fucked up")
			return
		}

		if userHasSubmittedRatingQuery.Next() {
			var count int64

			scanErr := userHasSubmittedRatingQuery.Scan(&count)

			userHasSubmittedRatingQuery.Close()

			if scanErr != nil {
				ctx.String(http.StatusOK, "because server fucked up")
				return
			}

			if count != 0 {
				ctx.String(http.StatusOK, fmt.Sprintf("alreadyvoted\n%.2f", totalRating))
				return
			}
		}

		ctx.String(http.StatusOK, fmt.Sprintf("ok\n%.2f", totalRating))
		return
	}
}
