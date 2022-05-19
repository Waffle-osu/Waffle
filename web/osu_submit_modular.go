package web

import (
	"Waffle/database"
	"Waffle/helpers"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type ScoreSubmission struct {
	FileHash            string
	Username            string
	OnlineScoreChecksum string
	Count300            int
	Count100            int
	Count50             int
	CountGeki           int
	CountKatu           int
	CountMiss           int
	TotalScore          int
	MaxCombo            int
	Perfect             bool
	Ranking             string
	EnabledMods         int
	Passed              bool
	Playmode            int
	Date                string
	ClientVersion       string
	ParsedSuccessfully  bool
}

func parseScoreString(score string) ScoreSubmission {
	splitScore := strings.Split(score, ":")

	count300, parseErr := strconv.Atoi(splitScore[3])
	count100, parseErr := strconv.Atoi(splitScore[4])
	count50, parseErr := strconv.Atoi(splitScore[5])
	countGeki, parseErr := strconv.Atoi(splitScore[6])
	countKatu, parseErr := strconv.Atoi(splitScore[7])
	countMiss, parseErr := strconv.Atoi(splitScore[8])
	totalScore, parseErr := strconv.Atoi(splitScore[9])
	maxCombo, parseErr := strconv.Atoi(splitScore[10])
	mods, parseErr := strconv.Atoi(splitScore[13])
	playmode, parseErr := strconv.Atoi(splitScore[15])

	if parseErr != nil {
		return ScoreSubmission{
			ParsedSuccessfully: false,
		}
	}

	perfect := false
	passed := false

	if splitScore[11] == "True" {
		perfect = true
	}

	if splitScore[14] == "True" {
		passed = true
	}

	scoreSubmission := ScoreSubmission{
		FileHash:            splitScore[0],
		Username:            splitScore[1],
		OnlineScoreChecksum: splitScore[2],
		Count300:            count300,
		Count100:            count100,
		Count50:             count50,
		CountGeki:           countGeki,
		CountKatu:           countKatu,
		CountMiss:           countMiss,
		TotalScore:          totalScore,
		MaxCombo:            maxCombo,
		Perfect:             perfect,
		Ranking:             splitScore[12],
		EnabledMods:         mods,
		Passed:              passed,
		Playmode:            playmode,
		Date:                splitScore[16],
		ClientVersion:       splitScore[17],
		ParsedSuccessfully:  true,
	}

	return scoreSubmission
}

func HandleOsuSubmit(ctx *gin.Context) {
	score := ctx.PostForm("score")
	password := ctx.PostForm("pass")
	wasExit := ctx.PostForm("x")
	failTime := ctx.PostForm("ft")
	clientHash := ctx.PostForm("s")
	//processList := ctx.PostForm("pl")

	//validate that parameters have indeed been sent
	if score == "" || password == "" || clientHash == "" {
		ctx.String(http.StatusBadRequest, "error: bad score submission")
		return
	}

	//peppy's score submission returns a key value pair with information about the beatmap and ranking and score changes
	//formatted like this: "key:value|key:value|key:value"
	//chartName:Overall Ranking|chartId:overall|toNextRank:123

	//peppy's score submission back then has these keys:
	//beatmapId            :: Beatmap ID
	//beatmapSetId         :: Beatmap Set ID
	//beatmapPlaycount     :: Beatmap Playcount
	//beatmapPasscount     :: Beatmap Passcount
	//approvedDate         :: When the Map was Approved
	//chartId              :: ID of a Chart, if it's just a normal score submission that goes to the main ranking, write "Overall Ranking"
	//chartName            :: Name of the Chart, if it's just a normal score submission that goes to the main ranking, write "overall"
	//chartEndDate         :: End Date of the Chart, leave empty if it's just a normal score submission
	//beatmapRankingBefore :: User's old rank on the beatmap
	//beatmapRankingAfter  :: User's rank on the beatmap now
	//rankedScoreBefore    :: User's old ranked score
	//rankedScoreAfter     :: User's ranked score now
	//totalScoreBefore     :: User's old total score
	//totalScoreAfter      :: User's total score now
	//playCountBefore      :: User's old playcount score
	//accuracyAfter        :: User's accuracy now
	//accuracyBefore       :: User's old accuracy
	//rankBefore           :: User's old rank
	//rankAfter            :: User's rank now
	//toNextRank           :: How much score until next leaderboard spot on the beatmap
	//toNextRankUser       :: How much more ranked score until the next ranked leaderboard spot
	//achievements         :: all achieved achievements in that play

	//alternatively, if an error were to occur, you return "error: what kind of error happened" the space after the : is important
	//there are some errors that the client itself will display an error for, these are:
	//"error: nouser"   :: For when the User doesn't exist
	//"error: pass"     :: For when the User's password is incorrect
	//"error: inactive" :: For when the User's account isn't activated
	//"error: ban"      :: For when the User is banned
	//"error: beatmap"  :: For when the beatmap is not available for ranking
	//"error: disabled" :: For when the Mode/Mod is currently disabled for ranking
	//"error: oldver"   :: For when the User's client is too old to submit scores

	scoreSubmissionResponse := make(map[string]string)

	//We don't have charts yet, so we just submit to the overall ranking
	scoreSubmissionResponse["chartName"] = "Overall Ranking"
	scoreSubmissionResponse["chartId"] = "overall"
	scoreSubmissionResponse["chartEndDate"] = ""

	scoreSubmission := parseScoreString(score)

	//fail the submission if the score wasnt parsed right
	if scoreSubmission.ParsedSuccessfully != true {
		ctx.String(http.StatusBadRequest, "error: bad score submission")
		return
	}

	userId, authSuccess := database.AuthenticateUser(scoreSubmission.Username, password)

	//server failure
	if userId == -2 {
		ctx.String(http.StatusOK, "error: fetch fail")
		return
	}

	//user not found
	if userId == -1 {
		ctx.String(http.StatusOK, "error: nouser")
		return
	}

	//wrong password
	if authSuccess == false {
		ctx.String(http.StatusOK, "error: pass")
		return
	}

	stringPerfect := "False"

	if scoreSubmission.Perfect == true {
		stringPerfect = "True"
	}

	stringPassed := "False"

	if scoreSubmission.Passed == true {
		stringPassed = "True"
	}

	//validate onlinescorechecksum
	onlineScoreChecksumInput := fmt.Sprintf("%do14%d%ds%d%duu%s%d%s%s%d%s%dQ%s%d%s%s%s",
		scoreSubmission.Count100+scoreSubmission.Count300,
		scoreSubmission.Count50,
		scoreSubmission.CountGeki,
		scoreSubmission.CountKatu,
		scoreSubmission.CountMiss,
		scoreSubmission.FileHash,
		scoreSubmission.MaxCombo,
		stringPerfect,
		scoreSubmission.Username,
		scoreSubmission.TotalScore,
		scoreSubmission.Ranking,
		scoreSubmission.EnabledMods,
		stringPassed,
		scoreSubmission.Playmode,
		scoreSubmission.ClientVersion,
		scoreSubmission.Date,
		clientHash)

	onlineScoreChecksumHashed := md5.Sum([]byte(onlineScoreChecksumInput))
	onlineScoreChecksumHashedString := hex.EncodeToString(onlineScoreChecksumHashed[:])

	if scoreSubmission.OnlineScoreChecksum != onlineScoreChecksumHashedString {
		ctx.String(http.StatusOK, "error: invalid score")
		return
	}

	//get users stats
	userFetchResult, userStats := database.UserStatsFromDatabase(uint64(userId), int8(scoreSubmission.Playmode))

	if userFetchResult != 0 {
		ctx.String(http.StatusOK, "error: nouser")
		return
	}

	helpers.Logger.Printf("[Web@ScoreSubmit] Got Score Submission from ID: %d; wasExit: %s; failTime: %s; clientHash: %s", userId, wasExit, failTime, clientHash)

	//save old values
	scoreSubmissionResponse["rankedScoreBefore"] = strconv.FormatUint(userStats.RankedScore, 10)
	scoreSubmissionResponse["totalScoreBefore"] = strconv.FormatUint(userStats.TotalScore, 10)
	scoreSubmissionResponse["playCountBefore"] = strconv.FormatUint(userStats.Playcount, 10)
	scoreSubmissionResponse["accuracyBefore"] = strconv.FormatFloat(float64(userStats.Accuracy), 'f', 2, 64)
	scoreSubmissionResponse["rankBefore"] = strconv.FormatUint(userStats.Rank, 10)

	//get map via the filehash
	beatmapFetchResult, scoreBeatmap := database.BeatmapsGetByMd5(scoreSubmission.FileHash)

	if beatmapFetchResult != 0 {
		ctx.String(http.StatusOK, "error: beatmap")
		return
	}

	//check for pending or unsubmitted status
	if scoreBeatmap.RankingStatus == database.BeatmapsDatabaseStatusPending || scoreBeatmap.RankingStatus == database.BeatmapsDatabaseStatusUnsubmitted {
		ctx.String(http.StatusOK, "error: beatmap")
		return
	}

	//Check for duplicate score
	duplicateScoreCheckQuery, duplicateScoreCheckQueryErr := database.Database.Query("SELECT COUNT(*) AS 'count' FROM waffle.scores WHERE score_hash = ?", scoreSubmission.OnlineScoreChecksum)

	if duplicateScoreCheckQueryErr != nil {
		ctx.String(http.StatusOK, "error: server error")

		if duplicateScoreCheckQuery != nil {
			duplicateScoreCheckQuery.Close()
		}

		return
	}

	if duplicateScoreCheckQuery.Next() {
		var count int64

		scanErr := duplicateScoreCheckQuery.Scan(&count)

		if scanErr != nil {
			ctx.String(http.StatusOK, "error: server error")

			if duplicateScoreCheckQuery != nil {
				duplicateScoreCheckQuery.Close()
			}
		}

		if count != 0 {
			ctx.String(http.StatusOK, "error: no duplicate scores!")

			if duplicateScoreCheckQuery != nil {
				duplicateScoreCheckQuery.Close()
			}

			return
		}
	} else {
		ctx.String(http.StatusOK, "error: server error")

		if duplicateScoreCheckQuery != nil {
			duplicateScoreCheckQuery.Close()
		}

		return
	}

	if duplicateScoreCheckQuery != nil {
		duplicateScoreCheckQuery.Close()
	}

	//save beatmap information
	scoreSubmissionResponse["beatmapId"] = strconv.FormatInt(int64(scoreBeatmap.BeatmapId), 10)
	scoreSubmissionResponse["beatmapsetId"] = strconv.FormatInt(int64(scoreBeatmap.BeatmapsetId), 10)
	scoreSubmissionResponse["approvedDate"] = scoreBeatmap.ApproveDate

	//query for play and passcount
	passPlayCountsQuery, passPlayCountsQueryErr := database.Database.Query("SELECT x.playcount, y.passcount FROM (SELECT COUNT(*) AS 'playcount' FROM waffle.scores WHERE beatmap_id = ? AND playmode = ?) AS x, (SELECT COUNT(*) AS 'passcount' FROM waffle.scores WHERE beatmap_id = ? AND playmode = ? AND passed = 1) AS y", scoreBeatmap.BeatmapId, int8(scoreSubmission.Playmode), scoreBeatmap.BeatmapId, int8(scoreSubmission.Playmode))

	//if we ever error, just send back 0
	if passPlayCountsQueryErr != nil {
		scoreSubmissionResponse["beatmapPlaycount"] = "1"
		scoreSubmissionResponse["beatmapPasscount"] = "0"
	} else {
		var playcount, passcount int64

		if passPlayCountsQuery.Next() {
			passPlayCountsQuery.Scan(&playcount, &passcount)

			scoreSubmissionResponse["beatmapPlaycount"] = strconv.FormatInt(playcount+1, 10)
			scoreSubmissionResponse["beatmapPasscount"] = strconv.FormatInt(passcount, 10)
		} else {
			scoreSubmissionResponse["beatmapPlaycount"] = "1"
			scoreSubmissionResponse["beatmapPasscount"] = "0"
		}
	}

	//get users best score
	scoreQueryResult, bestLeaderboardScore, _, _ := database.ScoresGetUserLeaderboardBest(scoreBeatmap.BeatmapId, uint64(userId), int8(scoreSubmission.Playmode))
	bestLeaderboardScoreExists := 0

	if scoreQueryResult == -2 {
		ctx.String(http.StatusOK, "error: server error")
		return
	}

	if scoreQueryResult == 0 {
		bestLeaderboardScoreExists = 1
	}

	//Increase playcount by 1
	userStats.Playcount++

	oldLeaderboardPlace := int64(0)

	if (bestLeaderboardScoreExists == 1 && bestLeaderboardScore.Score < scoreSubmission.TotalScore) || ((bestLeaderboardScore.Passed == 0 && scoreQueryResult == 0) && scoreSubmission.Passed == true) {
		queryResult, oldLeaderboardPlaceResult := database.ScoresGetBeatmapLeaderboardPlace(bestLeaderboardScore.ScoreId, int32(bestLeaderboardScore.BeatmapId))

		if queryResult != 0 {
			oldLeaderboardPlace = 0
		} else {
			oldLeaderboardPlace = oldLeaderboardPlaceResult
		}

		userStats.TotalScore -= uint64(bestLeaderboardScore.Score)
		userStats.Hit300 -= uint64(bestLeaderboardScore.Hit300)
		userStats.Hit100 -= uint64(bestLeaderboardScore.Hit100)
		userStats.Hit50 -= uint64(bestLeaderboardScore.Hit50)
		userStats.HitMiss -= uint64(bestLeaderboardScore.HitMiss)
		userStats.HitGeki -= uint64(bestLeaderboardScore.HitGeki)
		userStats.HitKatu -= uint64(bestLeaderboardScore.HitKatu)

		//Set that there is no best score anymore
		bestLeaderboardScoreExists = 0

		//Overwrite in database
		overwriteBestLeaderboardScoreQuery, overwriteBestLeaderboardScoreQueryErr := database.Database.Query("UPDATE waffle.scores SET leaderboard_best = 0 WHERE score_id = ?", bestLeaderboardScore.ScoreId)

		if overwriteBestLeaderboardScoreQueryErr != nil {
			ctx.String(http.StatusInternalServerError, "error: server error")
			return
		}

		if overwriteBestLeaderboardScoreQuery != nil {
			overwriteBestLeaderboardScoreQuery.Close()
		}
	}

	scoreSubmissionResponse["beatmapRankingBefore"] = strconv.FormatInt(oldLeaderboardPlace, 10)

	if bestLeaderboardScoreExists == 0 {
		userStats.TotalScore += uint64(scoreSubmission.TotalScore)
		userStats.Hit300 += uint64(scoreSubmission.Count300)
		userStats.Hit100 += uint64(scoreSubmission.Count100)
		userStats.Hit50 += uint64(scoreSubmission.Count50)
		userStats.HitMiss += uint64(scoreSubmission.CountMiss)
		userStats.HitGeki += uint64(scoreSubmission.CountGeki)
		userStats.HitKatu += uint64(scoreSubmission.CountKatu)

		userStats.Level = float64(helpers.GetLevelFromScore(userStats.TotalScore))
	}

	switch scoreSubmission.Playmode {
	case 0:
		userStats.Accuracy = helpers.CalculateGlobalAccuracyOsu(userStats.Hit50, userStats.Hit100, userStats.Hit300, userStats.HitGeki, userStats.HitKatu, userStats.HitMiss)
		break
	case 1:
		userStats.Accuracy = helpers.CalculateGlobalAccuracyTaiko(userStats.Hit50, userStats.Hit100, userStats.Hit300, userStats.HitGeki, userStats.HitKatu, userStats.HitMiss)
		break
	case 2:
		userStats.Accuracy = helpers.CalculateGlobalAccuracyCatch(userStats.Hit50, userStats.Hit100, userStats.Hit300, userStats.HitGeki, userStats.HitKatu, userStats.HitMiss)
		break
	}

	queryPerfect := int8(0)
	queryPassed := int8(0)
	queryLeaderboardBest := int8(0)
	queryMapsetBest := int8(0)

	mapsetBestScoreQueryResult, mapsetBestScore := database.ScoresGetBeatmapsetBestUserScore(scoreBeatmap.BeatmapsetId, uint64(userId), int8(scoreSubmission.Playmode))
	bestMapsetScoreExists := 0

	if mapsetBestScoreQueryResult == 0 {
		bestMapsetScoreExists = 1
	}

	if bestMapsetScoreExists == 1 && mapsetBestScore.Score < scoreSubmission.TotalScore && scoreSubmission.Passed && scoreBeatmap.RankingStatus != 2 {
		//I like to do this in 2 steps, makes me feel better
		userStats.RankedScore -= uint64(mapsetBestScore.Score)

		bestMapsetScoreExists = 0

		//Overwrite in database
		overwriteBestMapsetScoreQuery, overwriteBestMapsetScoreQueryErr := database.Database.Query("UPDATE waffle.scores SET mapset_best = 0 WHERE score_id = ?", mapsetBestScore.ScoreId)

		if overwriteBestMapsetScoreQueryErr != nil {
			ctx.String(http.StatusInternalServerError, "error: server error")
			return
		}

		if overwriteBestMapsetScoreQuery != nil {
			overwriteBestMapsetScoreQuery.Close()
		}
	}

	if bestMapsetScoreExists == 0 && scoreSubmission.Passed && scoreBeatmap.RankingStatus != 2 {
		userStats.RankedScore += uint64(scoreSubmission.TotalScore)
	}

	if bestLeaderboardScoreExists == 1 {
		queryLeaderboardBest = 0
	} else {
		queryLeaderboardBest = 1
	}

	if bestMapsetScoreExists == 1 {
		queryMapsetBest = 0
	} else {
		if scoreSubmission.Passed && scoreBeatmap.RankingStatus != 2 {
			queryMapsetBest = 1
		}
	}

	if scoreSubmission.Passed {
		queryPassed = 1
	}

	if scoreSubmission.Perfect {
		queryPerfect = 1
	}

	insertScoreQuery, insertScoreQueryErr := database.Database.Query("INSERT INTO waffle.scores (beatmap_id, beatmapset_id, user_id, playmode, score, max_combo, ranking, hit300, hit100, hit50, hitMiss, hitGeki, hitKatu, enabled_mods, perfect, passed, leaderboard_best, mapset_best, score_hash) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", scoreBeatmap.BeatmapId, scoreBeatmap.BeatmapsetId, userId, int8(scoreSubmission.Playmode), scoreSubmission.TotalScore, scoreSubmission.MaxCombo, scoreSubmission.Ranking, scoreSubmission.Count300, scoreSubmission.Count100, scoreSubmission.Count50, scoreSubmission.CountMiss, scoreSubmission.CountGeki, scoreSubmission.CountKatu, scoreSubmission.EnabledMods, queryPerfect, queryPassed, queryLeaderboardBest, queryMapsetBest, scoreSubmission.OnlineScoreChecksum)

	if insertScoreQueryErr != nil {
		ctx.String(http.StatusInternalServerError, "error: server error")
		return
	}

	if insertScoreQuery != nil {
		insertScoreQuery.Close()
	}

	scoreSubmissionResponse["rankedScoreAfter"] = strconv.FormatUint(userStats.RankedScore, 10)
	scoreSubmissionResponse["totalScoreAfter"] = strconv.FormatUint(userStats.TotalScore, 10)
	scoreSubmissionResponse["playCountAfter"] = strconv.FormatUint(userStats.Playcount, 10)
	scoreSubmissionResponse["accuracyAfter"] = strconv.FormatFloat(float64(userStats.Accuracy), 'f', 2, 64)

	updateUserStatsQuery, updateUserStatsQueryErr := database.Database.Query("UPDATE waffle.stats SET ranked_score = ?, total_score = ?, hit300 = ?, hit100 = ?, hit50 = ?, hitMiss = ?, hitGeki = ?, hitKatu = ?, user_level = ?, playcount = ?, accuracy = ? WHERE user_id = ? AND mode = ?", userStats.RankedScore, userStats.TotalScore, userStats.Hit300, userStats.Hit100, userStats.Hit50, userStats.HitMiss, userStats.HitGeki, userStats.HitKatu, userStats.Level, userStats.Playcount, userStats.Accuracy, userId, int8(scoreSubmission.Playmode))

	if updateUserStatsQueryErr != nil {
		ctx.String(http.StatusInternalServerError, "error: server error")
		return
	}

	if updateUserStatsQuery != nil {
		insertScoreQuery.Close()
	}

	newRankQuery, newRankQueryErr := database.Database.Query("SELECT `rank` FROM (SELECT user_id, mode, ROW_NUMBER() OVER (ORDER BY ranked_score DESC) AS 'rank' FROM waffle.stats WHERE mode = ? AND user_id != 1) t WHERE user_id = ?", int8(scoreSubmission.Playmode), userId)

	if newRankQueryErr != nil {
		ctx.String(http.StatusOK, "error: server error")
		return
	}

	var newRank int64

	if newRankQuery.Next() {
		scanErr := newRankQuery.Scan(&newRank)

		if scanErr != nil {
			ctx.String(http.StatusOK, "error: server error")
			return
		}

		newRankQuery.Close()

		scoreSubmissionResponse["rankAfter"] = strconv.FormatInt(newRank, 10)
	} else {
		//how tf would we get that far if the user wasnt there
		ctx.String(http.StatusOK, "error: nouser")

		if newRankQuery != nil {
			newRankQuery.Close()
		}

		return
	}

	if newRankQuery != nil {
		newRankQuery.Close()
	}

	newScoreIdGetQuery, newScoreIdGetQueryErr := database.Database.Query("SELECT score_id FROM (SELECT score_id, score_hash FROM waffle.scores WHERE score_hash = ?) t", scoreSubmission.OnlineScoreChecksum)
	newScoreId := int64(-1)

	if newScoreIdGetQueryErr != nil {
		ctx.String(http.StatusOK, "error: server error")

		if newScoreIdGetQuery != nil {
			newScoreIdGetQuery.Close()
		}

		return
	}

	if newScoreIdGetQuery.Next() {
		var scoreId int64

		scanErr := newScoreIdGetQuery.Scan(&scoreId)

		if scanErr != nil {
			ctx.String(http.StatusOK, "error: server error")
			return
		}

		newScoreIdGetQuery.Close()

		newScoreId = scoreId
	} else {
		//how tf would we get that far if the user wasnt there
		ctx.String(http.StatusOK, "error: nouser")

		if newScoreIdGetQuery != nil {
			newScoreIdGetQuery.Close()
		}

		return
	}

	if newScoreIdGetQuery != nil {
		newScoreIdGetQuery.Close()
	}

	newLeaderboardRankQueryResult, newLeaderboardRank := database.ScoresGetBeatmapLeaderboardPlace(uint64(newScoreId), scoreBeatmap.BeatmapId)

	if newLeaderboardRankQueryResult != 0 {
		newLeaderboardRank = 0
	}

	scoreSubmissionResponse["beatmapRankingAfter"] = strconv.FormatInt(newLeaderboardRank, 10)

	//If the user isn't rank 1, get how much score they need for the next rank
	if newRank != 1 {
		nextRankScoreQuery, nextRankScoreQueryErr := database.Database.Query("SELECT * FROM (SELECT users.username, stats.user_id, stats.ranked_score, stats.mode, ROW_NUMBER() OVER (ORDER BY ranked_score DESC) AS 'rank' FROM waffle.stats LEFT JOIN users ON stats.user_id = users.user_id WHERE mode = ?) t WHERE `rank` = ?", int8(scoreSubmission.Playmode), userStats.Rank-1)

		if nextRankScoreQueryErr != nil {
			ctx.String(http.StatusOK, "error: server error")
			return
		}

		if nextRankScoreQuery.Next() {
			var username string
			partUserStats := database.UserStats{}

			scanErr := nextRankScoreQuery.Scan(&username, &partUserStats.UserID, &partUserStats.RankedScore, &partUserStats.Mode, &partUserStats.Rank)

			if scanErr != nil {
				ctx.String(http.StatusOK, "error: server error")
				return
			}

			nextRankScoreQuery.Close()

			scoreSubmissionResponse["toNextRank"] = strconv.FormatInt(int64(partUserStats.RankedScore-userStats.RankedScore), 10)
			scoreSubmissionResponse["toNextRankUser"] = username
		} else {
			//how tf would we get that far if the user wasnt there
			ctx.String(http.StatusOK, "error: nouser")

			if nextRankScoreQuery != nil {
				nextRankScoreQuery.Close()
			}

			return
		}

		if nextRankScoreQuery != nil {
			nextRankScoreQuery.Close()
		}
	} else {
		scoreSubmissionResponse["toNextRank"] = "0"
		scoreSubmissionResponse["toNextRankUser"] = ""
	}

	//save failtime
	failTimeParsed, failTimeParseErr := strconv.ParseInt(failTime, 10, 64)

	if failTimeParseErr == nil {
		wasExitParsed := int8(0)

		if wasExit == "1" {
			wasExitParsed = 1
		}

		insertFailTimeQuery, _ := database.Database.Query("INSERT INTO waffle.failtimes (failtime, beatmap_id, score_id, was_exit) VALUES (?, ?, ?, ?)", failTimeParsed, scoreBeatmap.BeatmapId, newScoreId, wasExitParsed)

		if insertFailTimeQuery != nil {
			insertFailTimeQuery.Close()
		}
	}

	returnString := ""

	returnString += "beatmapId:" + scoreSubmissionResponse["beatmapId"] + "|"
	returnString += "beatmapSetId:" + scoreSubmissionResponse["beatmapSetId"] + "|"
	returnString += "beatmapPlaycount:" + scoreSubmissionResponse["beatmapPlaycount"] + "|"
	returnString += "beatmapPasscount:" + scoreSubmissionResponse["beatmapPasscount"] + "|"
	returnString += "approvedDate:" + scoreSubmissionResponse["approvedDate"]

	returnString += "\n"

	returnString += "chartId:" + scoreSubmissionResponse["chartId"] + "|"
	returnString += "chartName:" + scoreSubmissionResponse["chartName"] + "|"
	returnString += "chartEndDate:" + scoreSubmissionResponse["chartEndDate"] + "|"
	returnString += "beatmapRankingBefore:" + scoreSubmissionResponse["beatmapRankingBefore"] + "|"
	returnString += "beatmapRankingAfter:" + scoreSubmissionResponse["beatmapRankingAfter"] + "|"
	returnString += "rankedScoreBefore:" + scoreSubmissionResponse["rankedScoreBefore"] + "|"
	returnString += "rankedScoreAfter:" + scoreSubmissionResponse["rankedScoreAfter"] + "|"
	returnString += "totalScoreBefore:" + scoreSubmissionResponse["totalScoreBefore"] + "|"
	returnString += "totalScoreAfter:" + scoreSubmissionResponse["totalScoreAfter"] + "|"
	returnString += "playCountBefore:" + scoreSubmissionResponse["playCountBefore"] + "|"
	returnString += "accuracyBefore:" + scoreSubmissionResponse["accuracyBefore"] + "|"
	returnString += "accuracyAfter:" + scoreSubmissionResponse["accuracyAfter"] + "|"
	returnString += "rankBefore:" + scoreSubmissionResponse["rankBefore"] + "|"
	returnString += "rankAfter:" + scoreSubmissionResponse["rankAfter"] + "|"
	returnString += "toNextRank:" + scoreSubmissionResponse["toNextRank"] + "|"
	returnString += "toNextRankUser:" + scoreSubmissionResponse["toNextRankUser"] + "|"
	returnString += "achievements:" + scoreSubmissionResponse["achievements"]

	//make sure there's no trailing |
	returnString = strings.TrimSuffix(returnString, "|")

	ctx.String(http.StatusOK, returnString+"\n")

	replay, replayGetErr := ctx.FormFile("score")

	if replayGetErr == nil {
		ctx.SaveUploadedFile(replay, fmt.Sprintf("replays/%d", newScoreId))
	}
}
