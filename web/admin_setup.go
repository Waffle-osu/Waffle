package web

import (
	_ "embed"
	"github.com/gin-gonic/gin"
	"net/http"
)

//go:embed admin_setup_updater.html
var setup string

func HandleUpdaterAdminSetup(ctx *gin.Context) {
	ctx.String(http.StatusOK, setup)
}
