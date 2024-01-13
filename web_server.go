package main

import (
	"Waffle/api"
	"Waffle/site"
	"Waffle/web"
	"bytes"
	"context"
	"fmt"

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

	//BSS b1815
	// ginServer.POST("/web/osu-bmsubmit-getid5.php", bss.HandleGetId5)

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

	//Site:
	ginServer.GET("/", func(ctx *gin.Context) {
		buffer := new(bytes.Buffer)
		err := site.Index(site.UsageStatistics{
			TotalUsers:  fmt.Sprintf("%d", 24),
			OnlineUsers: fmt.Sprintf("%d", 3),
		}).Render(context.Background(), buffer)

		if err != nil {
			ctx.String(500, err.Error())
			return
		}

		str := buffer.String()

		ctx.Header("Content-Type", "text/html")
		ctx.String(200, str)
	})

	ginServer.GET("/site-style.css", func(ctx *gin.Context) {
		ctx.String(200, `
			.container-wrapper {
				display: grid;
				grid-template-columns: 1fr;
				grid-template-rows: 100vh;
				align-items: center;
				justify-items: center;
				height: 99%;
			}

			.container {
				color: white;
				width: 75%;
				height: 90%;
				display: grid; 
				grid-template-columns: 1fr 1fr 1fr 1fr 1fr; 
				grid-template-rows: 0.75fr 1fr 1fr 1fr 1fr; 
				gap: 0px 0px; 
				grid-template-areas: 
				"side-nav-bar top-info-bar top-info-bar top-info-bar top-info-bar"
				"side-nav-bar content content content content"
				"side-nav-bar content content content content"
				"side-nav-bar content content content content"
				"side-nav-bar content content content content"; 
			}

			.side-nav-bar { 
				padding-left: 16px;
				padding-top: 8px;
				background-color: rgb(50,50,50);
				grid-area: side-nav-bar; 
				margin-right: 24px;
			}

			.top-info-bar { 
				padding-left: 16px;
				padding-top: 16px;
				background-color: rgb(50,50,50);
				grid-area: top-info-bar; 
				margin-bottom: 24px;
			}

			.content {
				padding-left: 16px;
				background-color: rgb(50,50,50);
				grid-area: content; 
			}
		`)
	})

	ginServer.Run("127.0.0.1:80")
}
