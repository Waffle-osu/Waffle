package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RunOsuWeb() {
	ginServer := gin.Default()
	ginServer.SetTrustedProxies(nil)

	ginServer.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello, World!")
	})

	ginServer.POST("/web/osu-screenshot.php", HandleOsuScreenshot)

	ginServer.Run("127.0.0.1:80")
}
