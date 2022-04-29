package web

import (
	"Waffle/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleOsuUpdate2(ctx *gin.Context) {
	fileHash := database.GetOsuExecutableHash()

	if fileHash == "" {
		ctx.String(http.StatusInternalServerError, "")
		return
	}

	ctx.String(http.StatusOK, fileHash)
}
