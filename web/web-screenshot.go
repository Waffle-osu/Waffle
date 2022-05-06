package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func HandleOsuGetScreenshot(ctx *gin.Context) {
	//The osu! client opens the browser at /ss/:filename
	//where :filename the filename is, that you return in /web/osu-screenshot.php
	filename := ctx.Param("filename")

	screenshotBytes, readErr := os.ReadFile("screenshots/" + filename)

	if readErr != nil {
		ctx.String(http.StatusInternalServerError, "Failed to load screenshot!")
	}

	ctx.Data(http.StatusOK, "image/jpeg", screenshotBytes)
}
