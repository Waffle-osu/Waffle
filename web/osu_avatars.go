package web

import (
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

	ctx.File("avatars/" + filename + ".png")
}
