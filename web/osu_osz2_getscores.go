package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleOsuOsz2GetScores(ctx *gin.Context) {
	//Currently, has to still be done,
	//but the first line has to have the Ranking Status seperated with a | and whether the server has a .osz2 file of the map
	ctx.String(http.StatusOK, "2|false")
}
