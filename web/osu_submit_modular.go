package web

import (
	"Waffle/database"
	"Waffle/helpers"
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

	if splitScore[11] == "1" {
		perfect = true
	}

	if splitScore[14] == "1" {
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
	//replay, replayGetErr := ctx.FormFile("score")
	score := ctx.PostForm("score")
	password := ctx.PostForm("pass")
	wasExit := ctx.PostForm("x")
	failTime := ctx.PostForm("ft")
	clientHash := ctx.PostForm("s")
	processList := ctx.PostForm("pl")

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
	//rankScoreAfter       :: User's old rank
	//rankScoreAfter       :: User's rank now
	//toNextRank           :: How much score until next leaderboard spot on the beatmap
	//toNextRankUser       :: How much more ranked score until the next ranked leaderboard spot
	//achievements         :: all achieved achievements in that play

	//alternatively, if an error were to occur, you return "error: what kind of error happened" the space after the : is important
	//there are some errors that the client itself will display an error for, these are:
	//"error: nouser" :: For when the User doesn't exist
	//"error: pass" :: For when the User's password is incorrect
	//"error: inactive" :: For when the User's account isn't activated
	//"error: ban" :: For when the User is banned
	//"error: beatmap" :: For when the beatmap is not available for ranking
	//"error: disabled" :: For when the Mode/Mod is currently disabled for ranking
	//"error: oldver" :: For when the User's client is too old to submit scores

	chartInfo := make(map[string]string)

	chartInfo["chartName"] = "Overall Ranking"
	chartInfo["chartId"] = "overall"
	chartInfo["chartEndDate"] = ""

	scoreSubmission := parseScoreString(score)

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

	fetchResult, userStats := database.UserStatsFromDatabase(uint64(userId), int8(scoreSubmission.Playmode))

	if fetchResult != 0 {
		ctx.String(http.StatusOK, "error: nouser")
		return
	}

	helpers.Logger.Printf("[Web@ScoreSubmit] Got Score Submission from ID: %d; wasExit: %s; failTime: %s; clientHash: %s, processList: %s", userId, wasExit, failTime, clientHash, processList)

	chartInfo["rankedScoreBefore"] = strconv.FormatUint(userStats.RankedScore, 10)
	chartInfo["totalScoreBefore"] = strconv.FormatUint(userStats.TotalScore, 10)
	chartInfo["playCountBefore"] = strconv.FormatUint(userStats.Playcount, 10)
	chartInfo["accuracyBefore"] = strconv.FormatFloat(float64(userStats.Accuracy), 'f', 2, 64)
	chartInfo["rankBefore"] = strconv.FormatUint(userStats.Rank, 10)

	returnString := ""

	//Write out submission in the format the client expects
	for key, value := range chartInfo {
		returnString += key + ":" + value + "|"
	}

	//make sure there's no trailing |
	returnString = strings.TrimSuffix(returnString, "|")

	ctx.String(http.StatusOK, returnString+"\n")
}
