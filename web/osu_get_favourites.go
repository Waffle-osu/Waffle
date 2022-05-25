package web

import (
	"Waffle/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HandleOsuGetFavourites(ctx *gin.Context) {
	queryUsername := ctx.Query("u")
	queryPassword := ctx.Query("h")

	userId, authResult := database.AuthenticateUser(queryUsername, queryPassword)

	if userId == -2 {
		ctx.String(http.StatusInternalServerError, "")
		return
	}

	if !authResult {
		ctx.String(http.StatusUnauthorized, "")
		return
	}

	queryResult, favourites := database.GetUserFavourites(uint64(userId))

	if queryResult != 0 {
		ctx.String(http.StatusOK, "")
	}

	returnString := ""

	for i := 0; i != len(favourites); i++ {
		returnString += strconv.FormatInt(int64(favourites[i].BeatmapSetId), 10) + "\n"
	}

	ctx.String(http.StatusOK, returnString)
}
