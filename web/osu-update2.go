package web

import (
	"Waffle/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleOsuUpdate2(ctx *gin.Context) {
	fileHashPublic := database.UpdaterHashFromFilename("osu!.exe")
	fileHashTest := database.UpdaterHashFromFilename("osu!test.exe")

	if fileHashPublic == "" {
		ctx.String(http.StatusInternalServerError, "")
		return
	}

	result := "osu!.exe " + fileHashPublic + "\n"

	if fileHashTest != "" {
		result += "osu!test.exe" + fileHashTest + "\n"
	}
	
	ctx.String(http.StatusOK, result)
}
