package web

import (
	"Waffle/database"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

func HandleOsuScreenshot(ctx *gin.Context) {
	//The osu! client sends along its credentials, and as a Form File parameter, a JPEG image containing the screenshot
	//For whatever reason, no matter what option you pick in settings, it'll always upload JPEGs, never PNGs
	userId, authResult := database.AuthenticateUser(ctx.Query("u"), ctx.Query("p"))

	if userId == -2 {
		ctx.String(http.StatusInternalServerError, "Failed to upload Screenshot!")
		return
	}

	if authResult == false {
		ctx.String(http.StatusUnauthorized, "Failed to upload Screenshot!")
		return
	}

	if database.HasReachedMaxScreenshotQuota(uint64(userId)) {
		ctx.String(http.StatusOK, "chill_out_man")
		return
	}

	filename := fmt.Sprintf("%x-%x", userId, time.Now().Unix())

	if database.InsertNewScreenshot(uint64(userId), filename) == false {
		ctx.String(http.StatusOK, "an_error_occured")
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

	os.WriteFile("screenshots/"+filename, fileBuffer, 0644)

	ctx.String(http.StatusOK, filename)
}
