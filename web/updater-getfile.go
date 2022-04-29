package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func HandleUpdaterGetFile(ctx *gin.Context) {
	filename := ctx.Param("filename")

	fileBytes, readErr := os.ReadFile("release/" + filename)

	if readErr != nil {
		ctx.String(http.StatusInternalServerError, "Failed to download file!")
		return
	}

	ctx.Data(http.StatusOK, "waffle/blob", fileBytes)
}
