package web

import (
	"Waffle/database"

	"github.com/gin-gonic/gin"
)

func HandleOsuLogin(ctx *gin.Context) {
	username := ctx.Query("username")
	password := ctx.Query("password")

	_, success := database.AuthenticateUser(username, password)

	if success {
		ctx.String(200, "1")
	} else {
		ctx.String(200, "0")
	}
}
