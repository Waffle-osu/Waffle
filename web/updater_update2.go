package web

import (
	"Waffle/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleUpdaterUpdate2(ctx *gin.Context) {
	//Here the Updater is asking the server for all the Updater Items available,
	//You have to return all available Updater Things formatted accordingly
	//Format specification is in item.FormatUpdaterItem() they are all seperated by a new line

	result, items := database.UpdaterGetUpdaterItems()

	if result == -1 {
		ctx.String(http.StatusInternalServerError, "")
		return
	}

	returnString := ""

	for _, item := range items {
		returnString += item.FormatUpdaterItem()
	}

	ctx.String(http.StatusOK, returnString)
}
