package web

import (
	"Waffle/bancho/database"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

func HandleOsuScreenshot(ctx *gin.Context) {
	userId, authResult := database.AuthenticateUser(ctx.Query("u"), ctx.Query("p"))

	if userId == -2 {
		ctx.String(http.StatusInternalServerError, "Failed to upload Screenshot!")
		return
	}

	if authResult == false {
		ctx.String(http.StatusUnauthorized, "Failed to upload Screenshot!")
		return
	}

	//peppy sends a screenshot using multipart forms, he calls the label for the screenshot 'ss'
	screenshot, formErr := ctx.FormFile("ss")

	if formErr != nil {
		ctx.String(http.StatusBadRequest, "Failed to upload Screenshot!")
		return
	}

	//Open File
	ssFile, ssFileErr := screenshot.Open()

	if ssFileErr != nil {
		ctx.String(http.StatusBadRequest, "Failed to upload Screenshot!")
		return
	}

	//Make a buffer large enough to fit and read in
	fileBuffer := make([]byte, screenshot.Size)
	ssFile.Read(fileBuffer)

	filename := fmt.Sprintf("%x-%x", userId, time.Now().Unix())

	os.WriteFile("screenshots/"+filename, fileBuffer, 0644)

	ctx.String(http.StatusOK, filename)
}
