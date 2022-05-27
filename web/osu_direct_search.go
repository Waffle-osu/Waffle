package web

import (
	"Waffle/database"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DirectDisplayModeRanked    = 0
	DirectDisplayModePending   = 2
	DirectDisplayModeAll       = 4
	DirectDisplayModeGraveyard = 5
)

type DirectBeatmapQuery struct {
	BeatmapsetId  int
	Artist        string
	Title         string
	Creator       string
	HasVideo      int8
	HasStoryboard int8
	RatingSum     int64
	Votes         int64
	RankingStatus int8
	ApproveDate   string
}

func HandleOsuDirectSearch(ctx *gin.Context) {
	queryUsername := ctx.Query("u")
	queryPassword := ctx.Query("h")
	queryDisplayMode := ctx.Query("r")
	queryQuery := ctx.Query("q")

	userId, authResult := database.AuthenticateUser(queryUsername, queryPassword)

	if userId == -2 {
		ctx.String(http.StatusInternalServerError, "-1\nServer Error")
		return
	}

	if !authResult {
		ctx.String(http.StatusUnauthorized, "-1\nUnauthorized")
		return
	}

	rankedStatuses := ""

	switch queryDisplayMode {
	case "0":
		rankedStatuses = "1,2"
	case "2":
		rankedStatuses = "0,-1"
	case "5":
		rankedStatuses = "0,-1"
	case "4":
		rankedStatuses = "-1,0,1,2"
	}

	var beatmapRows *sql.Rows

	if queryQuery == "Newest" {
		newestSqlQuery := `
SELECT 
	result.beatmapset_id, 
	result.artist, 
	result.title, 		
	result.creator, 
	result.has_video, 
	result.has_storyboard, 
	result.final_rating_sum AS 'rating_sum', 
	result.final_votes AS 'votes',  
	result.ranking_status,
	result.approve_date
FROM (
	SELECT 
		beatmapsets.beatmapset_id, 
		beatmapsets.artist, 
		beatmapsets.title, 
		beatmapsets.creator, 
		beatmapsets.has_video, 
		beatmapsets.has_storyboard, 
		beatmap_ratings.rating_sum, 
		beatmap_ratings.votes, 
		beatmaps.approve_date,
		beatmaps.ranking_status,
		CASE WHEN beatmap_ratings.rating_sum IS NULL THEN 0 ELSE beatmap_ratings.rating_sum END AS 'final_rating_sum',
		CASE WHEN beatmap_ratings.votes IS NULL THEN 1 ELSE beatmap_ratings.votes END AS 'final_votes'
	FROM waffle.beatmapsets 
		LEFT JOIN waffle.beatmaps ON beatmaps.beatmapset_id = beatmapsets.beatmapset_id
		LEFT JOIN waffle.beatmap_ratings ON beatmap_ratings.beatmapset_id = beatmapsets.beatmapset_id
	WHERE ranking_status IN (%s)
	GROUP BY beatmapsets.beatmapset_id 
	ORDER BY approve_date DESC
	LIMIT 250
) result
		`

		newestQuery, newestQueryErr := database.Database.Query(fmt.Sprintf(newestSqlQuery, rankedStatuses))

		if newestQueryErr != nil {
			ctx.String(http.StatusOK, "-1\nNewest Query Failed...\n")
			return
		}

		beatmapRows = newestQuery
	} else if queryQuery == "Top Rated" {
		topRatedQuerySql := `
SELECT 
	result.beatmapset_id, 
	result.artist, 
	result.title, 		
	result.creator, 
	result.has_video, 
	result.has_storyboard, 
	result.final_rating_sum AS 'rating_sum', 
	result.final_votes AS 'votes',  
	result.ranking_status,
	result.approve_date
FROM (
	SELECT 
		beatmapsets.beatmapset_id, 
		beatmapsets.artist, 
		beatmapsets.title, 
		beatmapsets.creator, 
		beatmapsets.has_video, 
		beatmapsets.has_storyboard, 
		beatmap_ratings.rating_sum, 
		beatmap_ratings.votes, 
		beatmaps.approve_date,
		beatmaps.ranking_status,
		CASE WHEN beatmap_ratings.rating_sum IS NULL THEN 0 ELSE beatmap_ratings.rating_sum END AS 'final_rating_sum',
		CASE WHEN beatmap_ratings.votes IS NULL THEN 1 ELSE beatmap_ratings.votes END AS 'final_votes'
	FROM waffle.beatmapsets 
		LEFT JOIN waffle.beatmaps ON beatmaps.beatmapset_id = beatmapsets.beatmapset_id
		LEFT JOIN waffle.beatmap_ratings ON beatmap_ratings.beatmapset_id = beatmapsets.beatmapset_id
	WHERE ranking_status IN (%s)
	GROUP BY beatmapsets.beatmapset_id 
	ORDER BY (beatmap_ratings.rating_sum / (beatmap_ratings.votes + 1)) DESC
	LIMIT 250
) result
		`

		topRatedQuery, topRatedQueryErr := database.Database.Query(fmt.Sprintf(topRatedQuerySql, rankedStatuses))

		if topRatedQueryErr != nil {
			ctx.String(http.StatusOK, "-1\nTop Rated Query Failed...\n")
			return
		}

		beatmapRows = topRatedQuery
	} else {
		generalSearchSql := `
SELECT 
	result.beatmapset_id, 
	result.artist, 
	result.title, 
	result.creator, 
	result.has_video, 
	result.has_storyboard, 
	result.final_rating_sum AS 'rating_sum', 
	result.final_votes AS 'votes',  
	result.ranking_status,
	result.approve_date
FROM (
	SELECT 
		beatmapsets.beatmapset_id, 
		beatmapsets.artist, 
		beatmapsets.title, 
		beatmapsets.creator, 
		beatmapsets.has_video, 
		beatmapsets.has_storyboard, 
		beatmap_ratings.rating_sum, 
		beatmap_ratings.votes, 
		beatmaps.approve_date,
		beatmaps.ranking_status,
		CASE WHEN beatmap_ratings.rating_sum IS NULL THEN 0 ELSE beatmap_ratings.rating_sum END AS 'final_rating_sum',
		CASE WHEN beatmap_ratings.votes IS NULL THEN 1 ELSE beatmap_ratings.votes END AS 'final_votes'
	FROM waffle.beatmapsets 
		LEFT JOIN waffle.beatmaps ON beatmaps.beatmapset_id = beatmapsets.beatmapset_id
		LEFT JOIN waffle.beatmap_ratings ON beatmap_ratings.beatmapset_id = beatmapsets.beatmapset_id
	WHERE ranking_status IN (%s) AND (
		LOWER(title) LIKE CONCAT('%%', ?, '%%') OR 
		LOWER(artist) LIKE CONCAT('%%', ?, '%%') OR 
		LOWER(creator) LIKE CONCAT('%%', ?, '%%') OR
		LOWER(source) LIKE CONCAT('%%', ?, '%%') OR
		LOWER(tags) LIKE CONCAT('%%', ?, '%%')
	)
	GROUP BY beatmapsets.beatmapset_id 
	LIMIT 250
) result
		`

		formattedSql := fmt.Sprintf(generalSearchSql, rankedStatuses)

		searchQuery, searchQueryErr := database.Database.Query(formattedSql, queryQuery, queryQuery, queryQuery, queryQuery, queryQuery)

		if searchQueryErr != nil {
			ctx.String(http.StatusOK, "-1\nQuery Failed...\n")
			return
		}

		beatmapRows = searchQuery
	}

	returnString := "0\n"

	for beatmapRows.Next() {
		beatmap := DirectBeatmapQuery{}

		scanErr := beatmapRows.Scan(&beatmap.BeatmapsetId, &beatmap.Artist, &beatmap.Title, &beatmap.Creator, &beatmap.HasVideo, &beatmap.HasStoryboard, &beatmap.RatingSum, &beatmap.Votes, &beatmap.RankingStatus, &beatmap.ApproveDate)

		if scanErr != nil {
			beatmapRows.Close()
			ctx.String(http.StatusOK, "-1\nNewest Query Failed...\n")
			return
		}

		//Seperated by |
		//[0]:  Server Filename
		//[1]:  Artist
		//[2]:  Title
		//[3]:  Creator
		//[4]:  Ranked Status
		//[5]:  Rating
		//[6]:  Last Update
		//[7]:  Set ID
		//[8]:  Thread ID
		//[9]:  Has Video
		//[10]: Has Storyboard
		//[11]: File Size
		//[12]: File Size without Video

		if fileStats, err := os.Stat("oszs/" + strconv.FormatInt(int64(beatmap.BeatmapsetId), 10) + ".osz"); errors.Is(err, os.ErrNotExist) {

		} else {
			returnString += fmt.Sprintf("%s|%s|%s|%s|%d|%.2f|%s|%d|%d|%d|%d|%d|%d\n", strconv.FormatInt(int64(beatmap.BeatmapsetId), 10)+".osz", beatmap.Artist, beatmap.Title, beatmap.Creator, beatmap.RankingStatus, float64(beatmap.RatingSum)/float64(beatmap.Votes), beatmap.ApproveDate, beatmap.BeatmapsetId, 0, beatmap.HasVideo, beatmap.HasStoryboard, fileStats.Size(), fileStats.Size())
		}
	}

	beatmapRows.Close()

	ctx.String(http.StatusOK, returnString)
}
