package web

import (
	"Waffle/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleOsuUpdate2(ctx *gin.Context) {
	//Here the osu! client is asking for MD5 hashes of the client
	//Because there are 2 streams of the client, we're sending both of them
	//As the client will check for its own filename
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
