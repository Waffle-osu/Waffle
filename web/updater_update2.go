package web

import (
	"Waffle/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func HandleUpdaterUpdate2(ctx *gin.Context) {
	_, readErr := os.Stat(".env")

	if readErr != nil {
		ctx.String(http.StatusOK, "\n")
		return
	}

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
