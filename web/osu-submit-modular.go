package web

import (
	"Waffle/database"
	"Waffle/logger"
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

	scoreSubmission := parseScoreString(score)

	if scoreSubmission.ParsedSuccessfully != true {
		ctx.String(http.StatusBadRequest, "error: bad score submission")
		return
	}

	userId, authSuccess := database.AuthenticateUser(scoreSubmission.Username, password)

	if authSuccess == false {
		ctx.String(http.StatusOK, "error: pass")
		return
	}

	logger.Logger.Printf("[Web@ScoreSubmit] Got Score Submission from ID: %d; wasExit: %s; failTime: %s; clientHash: %s, processList: %s", userId, wasExit, failTime, clientHash, processList)
}
