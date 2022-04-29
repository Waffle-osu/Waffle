package web

import (
	"Waffle/database"
	_ "embed"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

//go:embed admin_post_success_updater.html
var postSuccess string

func HandleUpdaterAdminPostSettings(ctx *gin.Context) {
	formUsername := ctx.PostForm("mysql_username")
	formPassword := ctx.PostForm("mysql_password")
	formLocation := ctx.PostForm("mysql_location")
	formDatabase := ctx.PostForm("mysql_database")

	if formUsername == "" || formPassword == "" || formLocation == "" || formDatabase == "" {
		ctx.String(http.StatusBadRequest, "Actually fill in the form properly, thank you")
		return
	}

	envFile := "mysql_username=" + formUsername + "\n"
	envFile += "mysql_password=" + formPassword + "\n"
	envFile += "mysql_location=" + formLocation + "\n"
	envFile += "mysql_database=" + formDatabase + "\n"

	writeErr := os.WriteFile(".env", []byte(envFile), 0644)

	if writeErr != nil {
		ctx.String(http.StatusInternalServerError, ".env File failed to save successfully...")
		return
	}

	ctx.String(http.StatusOK, postSuccess)

	database.Initialize(formUsername, formPassword, formLocation, formDatabase)
}
