package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleUpdaterChangelog(ctx *gin.Context) {
	//The Updater here expects just a webpage with the changelog, go nuts here i guess,
	//Keep in mind the Updater is using Internet Explorer,

	ctx.String(http.StatusOK, "Insert Changelog here...")
}
