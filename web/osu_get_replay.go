package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func HandleGetReplay(ctx *gin.Context) {
	scoreId := ctx.Query("c")

	replay, error := os.ReadFile("replays/" + scoreId)
	fmt.Println("replay length: " + string(rune(len(replay))))

	if error != nil {
		ctx.Data(http.StatusOK, "image/png", []byte{})
		fmt.Println("replay get failed")
	} else {
		ctx.Data(http.StatusOK, "waffle/blob", replay)
	}
}
