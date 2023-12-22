package main

import (
	"Waffle/api"
	"Waffle/web"

	"github.com/gin-gonic/gin"
)

func RunWeb() {
	ginServer := gin.Default()
	ginServer.SetTrustedProxies(nil)

	// /weeb
	ginServer.POST("/web/osu-screenshot.php", web.HandleOsuScreenshot)
	ginServer.GET("/web/osu-title-image.php", web.HandleTitleImage)
	ginServer.POST("/web/osu-submit-modular.php", web.HandleOsuSubmit)
	ginServer.GET("/web/osu-osz2-getscores.php", web.HandleOsuGetLeaderboards)
	ginServer.GET("/web/osu-getscores6.php", web.HandleOsuGetLeaderboards)
	ginServer.GET("/web/osu-getreplay.php", web.HandleGetReplay)
	ginServer.GET("/web/osu-getfavourites.php", web.HandleOsuGetFavourites)
	ginServer.GET("/web/osu-addfavourite.php", web.HandleOsuAddFavourite)
	ginServer.POST("/web/osu-comment.php", web.HandleOsuComments)
	ginServer.GET("/rating/ingame-rate2.php", web.HandleOsuIngameRate2)
	ginServer.GET("/web/osu-search.php", web.HandleOsuDirectSearch)
	ginServer.GET("/web/maps/:filename", web.HandleOsuMapUpdate)

	// updater
	//ginServer.GET("/p/changelog", HandleUpdaterChangelog)
	//ginServer.GET("/release/update2.txt", HandleUpdaterUpdate2)
	//ginServer.GET("/release/update2.php", HandleOsuUpdate2)
	//ginServer.GET("/release/:filename", HandleUpdaterGetFile)

	//direct stuff
	ginServer.GET("/mt/:filename", web.HandleOsuGetDirectThumbnail)
	ginServer.GET("/mp3/preview/:filename", web.HandleOsuGetDirectMp3Preview)
	ginServer.GET("/d/:filename", web.HandleOsuDirectDownload)
	ginServer.GET("/web/osu-search-set.php", web.HandleDirectSearchSet)

	//avatars
	ginServer.GET("/a/:filename", web.HandleOsuGetAvatar)

	// screenshots
	ginServer.GET("/ss/:filename", web.HandleOsuGetScreenshot)

	//api
	ginServer.POST("/api/waffle-login", api.ApiHandleWaffleLogin)
	ginServer.POST("/api/waffle-site-register", api.ApiHandleWaffleRegister)

	//achievements
	ginServer.GET("/images/achievements/:filename", web.HandleOsuGetAchievementImage)

	ginServer.Run("127.0.0.1:80")
}
