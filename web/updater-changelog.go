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

	//The Updater here expects just a webpage with the changelog, go nuts here i guess,
	//Keep in mind the Updater is using Internet Explorer,

	ctx.String(http.StatusOK, "Insert Changelog here...")
}
