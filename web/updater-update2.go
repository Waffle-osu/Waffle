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

	result, items := database.GetUpdaterItems()

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
