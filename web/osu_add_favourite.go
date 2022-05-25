package web

import (
	"Waffle/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HandleOsuAddFavourite(ctx *gin.Context) {
	queryUsername := ctx.Query("u")
	queryPassword := ctx.Query("h")
	queryBeatmapSetId := ctx.Query("a")

	userId, authResult := database.AuthenticateUser(queryUsername, queryPassword)

	if userId == -2 {
		ctx.String(http.StatusInternalServerError, "Failed to Add Favourite!")
		return
	}

	if !authResult {
		ctx.String(http.StatusUnauthorized, "Failed to Add Favourite!")
		return
	}

	beatmapSetId, parseErr := strconv.ParseInt(queryBeatmapSetId, 10, 64)

	if parseErr != nil {
		ctx.String(http.StatusBadRequest, "Failed to Add Favourite!")
		return
	}

	database.FavouritesAddFavourite(uint64(userId), int32(beatmapSetId))

	ctx.String(http.StatusOK, "")
}
