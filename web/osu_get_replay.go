package web

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func HandleGetReplay(ctx *gin.Context) {
	scoreId := ctx.Query("c")

	replay, error := os.ReadFile("replays/" + scoreId)

	if error != nil {
		ctx.Data(http.StatusOK, "waffle/blob", []byte{})
		fmt.Println("replay get failed")
	} else {
		ctx.Data(http.StatusOK, "waffle/blob", replay)
	}
}
