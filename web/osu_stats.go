package web

import (
	"Waffle/database"
	"fmt"

	"github.com/gin-gonic/gin"
)

func HandleOsuStats(ctx *gin.Context) {
	username := ctx.Query("u")
	password := ctx.Query("p")

	userId, success := database.AuthenticateUser(username, password)

	if !success {
		ctx.String(200, "")

		return
	}

	querySuccess, stats := database.UserStatsFromDatabase(uint64(userId), 0)

	if querySuccess != 0 {
		ctx.String(200, "")

		return
	}

	returnString := fmt.Sprintf("%d|", stats.RankedScore)
	returnString += fmt.Sprintf("%.2f|", stats.Accuracy)
	returnString += "unused|"
	returnString += "unused|"
	returnString += fmt.Sprintf("%d|", stats.Rank)
	returnString += fmt.Sprintf("%d", userId)

	ctx.String(200, returnString)
}

func HandleOsuStatsOthers(ctx *gin.Context) {
	username := ctx.Query("u")
	// check := ctx.Query("c")

	queryResult, user := database.UserFromDatabaseByUsername(username)
	queryResult2, stats := database.UserStatsFromDatabase(user.UserID, 0)

	if queryResult != 0 || queryResult2 != 0 {
		ctx.String(200, "")

		return
	}

	returnString := fmt.Sprintf("%d|", stats.RankedScore)
	returnString += fmt.Sprintf("%.2f|", stats.Accuracy)
	returnString += "unused|"
	returnString += "unused|"
	returnString += fmt.Sprintf("%d|", stats.Rank)
	returnString += fmt.Sprintf("%d", user.UserID)

	ctx.String(200, returnString)
}
