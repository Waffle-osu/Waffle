package web

import (
	"errors"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func HandleOsuGetAchievementImage(ctx *gin.Context) {
	filename := ctx.Param("filename")

	if strings.Contains(filename, "\\") || strings.Contains(filename, "/") || strings.Contains(filename, "..") {
		ctx.String(200, "fuck off")
		return
	}

	if _, err := os.Stat("achievement_images/" + filename); errors.Is(err, os.ErrNotExist) {
		ctx.Data(404, "data/blob", []byte{})
		return
	}

	data, err := os.ReadFile("achievement_images/" + filename)

	if err != nil {
		ctx.Data(502, "data/blob", []byte{})
		return
	}

	ctx.Data(200, "data/blob", data)
}
