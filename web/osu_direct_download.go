package web

import (
	"Waffle/database"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func HandleOsuDirectDownload(ctx *gin.Context) {
	queryUsername := ctx.Query("u")
	queryPassword := ctx.Query("h")
	queryFilename := ctx.Param("filename")

	queryFilename = strings.TrimSuffix(queryFilename, "n")

	userId, authResult := database.AuthenticateUser(queryUsername, queryPassword)

	if userId == -2 {
		ctx.String(http.StatusInternalServerError, "")
		return
	}

	if !authResult {
		ctx.String(http.StatusUnauthorized, "")
		return
	}

	setIdFromFilename, parseErr := strconv.ParseInt(queryFilename, 10, 64)

	if parseErr != nil {
		ctx.String(http.StatusBadRequest, "")
		return
	}

	oszFilename := "oszs/" + queryFilename + ".osz"

	beatmapSetQueryErr, beatmapSet := database.BeatmapsetsGetBeatmapsetById(int32(setIdFromFilename))

	if beatmapSetQueryErr != 0 {
		ctx.String(http.StatusNotFound, "")
		return
	}

	returnFilename := fmt.Sprintf("%d %s - %s (%s).osz", beatmapSet.BeatmapsetId, beatmapSet.Artist, beatmapSet.Title, beatmapSet.Creator)

	if _, err := os.Stat(oszFilename); errors.Is(err, os.ErrNotExist) {
		ctx.String(http.StatusNotFound, "")
		return
	} else {
		ctx.FileAttachment(oszFilename, returnFilename)
	}
}
