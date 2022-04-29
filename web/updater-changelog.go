package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func HandleUpdaterChangelog(ctx *gin.Context) {
	_, readErr := os.Stat(".env")

	if readErr != nil {
		ctx.Redirect(http.StatusFound, "/admin/waffle_setup")
		return
	}

	ctx.String(http.StatusOK, "Insert Changelog here...")
}
