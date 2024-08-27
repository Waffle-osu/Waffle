package main

import (
	"Waffle/api"
	"Waffle/stream"
	"Waffle/web"
	"Waffle/web/bss"

	"github.com/gin-gonic/gin"
)

func RunWeb() {
	ginServer := gin.Default()
	ginServer.SetTrustedProxies(nil)

	// /web
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

	ginServer.GET("/web/osu-login.php", web.HandleOsuLogin)
	ginServer.GET("/web/osu-stat.php", web.HandleOsuStats)
	ginServer.GET("/web/osu-statoth.php", web.HandleOsuStatsOthers)

	//BSS b1815
	ginServer.POST("/web/osu-bmsubmit-post3.php", bss.HandlePost3)
	ginServer.POST("/web/osu-bmsubmit-getid5.php", bss.HandleGetId5)
	ginServer.POST("/web/osu-bmsubmit-upload.php", bss.HandleUpload)

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
	ginServer.GET("/forum/download.php", web.HandleOsuGetForumAvatar)

	// screenshots
	ginServer.GET("/ss/:filename", web.HandleOsuGetScreenshot)

	//api
	ginServer.POST("/api/waffle-login", api.ApiHandleWaffleLogin)
	ginServer.POST("/api/waffle-site-register", api.ApiHandleWaffleRegister)

	//achievements
	ginServer.GET("/images/achievements/:filename", web.HandleOsuGetAchievementImage)

	//osu!stream arcade thingy
	ginServer.POST("/stream/arcade-auth", stream.HandleArcadeAuth)
	ginServer.POST("/stream/arcade-link", stream.HandleArcadeLink)
	ginServer.POST("/stream/arcade-token", stream.HandleArcadeToken)

	ginServer.Run("127.0.0.1:80")
}
