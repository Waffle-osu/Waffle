package web

import (
	"errors"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func HandleOsuGetAvatar(ctx *gin.Context) {
	filename := ctx.Param("filename")

	if strings.HasSuffix(filename, "_000.png") {
		filename = strings.TrimSuffix(filename, "_000.png")
	} else {
		filename = strings.TrimSuffix(filename, "_000.jpg")
	}

	if _, err := os.Stat("avatars/" + filename + ".png"); errors.Is(err, os.ErrNotExist) {
		filename = "2"
	}

	ctx.File("avatars/" + filename + ".png")
}

func HandleOsuGetForumAvatar(ctx *gin.Context) {
	filename := ctx.Query("avatar")

	if _, err := os.Stat("avatars/" + filename + ".png"); errors.Is(err, os.ErrNotExist) {
		filename = "2"
	}

	ctx.File("avatars/" + filename + ".png")
}
