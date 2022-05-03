package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleOsuOsz2GetScores(ctx *gin.Context) {
	ctx.String(http.StatusOK, "2|false")
}
