package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleUpdaterChangelog(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Hello, World!")
}
