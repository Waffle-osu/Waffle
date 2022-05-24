package web

import (
	"github.com/gin-gonic/gin"
)

func RunOsuWeb() {
	ginServer := gin.Default()
	ginServer.SetTrustedProxies(nil)

	// /web
	ginServer.POST("/web/osu-screenshot.php", HandleOsuScreenshot)
	ginServer.GET("/web/osu-title-image.php", HandleTitleImage)
	ginServer.POST("/web/osu-submit-modular.php", HandleOsuSubmit)
	ginServer.GET("/web/osu-osz2-getscores.php", HandleOsuGetLeaderboards)
	ginServer.GET("/web/osu-getscores6.php", HandleOsuGetLeaderboards)
	ginServer.GET("/web/osu-getreplay.php", HandleGetReplay)

	// updater
	//ginServer.GET("/p/changelog", HandleUpdaterChangelog)
	//ginServer.GET("/release/update2.txt", HandleUpdaterUpdate2)
	//ginServer.GET("/release/update2.php", HandleOsuUpdate2)
	//ginServer.GET("/release/:filename", HandleUpdaterGetFile)

	// screenshots
	ginServer.GET("/ss/:filename", HandleOsuGetScreenshot)

	ginServer.Run("127.0.0.1:80")
}
