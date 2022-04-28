package web

import (
	"github.com/gin-gonic/gin"
)

func RunOsuWeb() {
	ginServer := gin.Default()
	ginServer.SetTrustedProxies(nil)

	ginServer.POST("/web/osu-screenshot.php", HandleOsuScreenshot)
	ginServer.GET("/web/osu-title-image.php", HandleTitleImage)
	ginServer.GET("/ss/:filename", HandleOsuGetScreenshot)

	ginServer.Run("127.0.0.1:80")
}
