package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func HandleUpdaterGetFile(ctx *gin.Context) {
	//The Updater just goes to /release/:filename
	//:filename is controlled by you, in the Updater Item Specification it's the Server Filename
	//You return just the binary file just like that
	filename := ctx.Param("filename")

	fileBytes, readErr := os.ReadFile("release/" + filename)

	if readErr != nil {
		ctx.String(http.StatusInternalServerError, "Failed to download file!")
		return
	}

	ctx.Data(http.StatusOK, "waffle/blob", fileBytes)
}
