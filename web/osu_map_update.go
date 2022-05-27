package web

import "github.com/gin-gonic/gin"

func HandleOsuMapUpdate(ctx *gin.Context) {
	filename := ctx.Param("filename")

	ctx.File("osus/" + filename)
}
