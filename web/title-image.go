package web

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func HandleTitleImage(ctx *gin.Context) {
	version := ctx.Query("v") //Client Edition, t = testers, p = public
	current := ctx.Query("c") //MD5 Hash of the client's currently stored title image
	clicked := ctx.Query("l") //Whether the image was clicked

	if clicked == "1" {
		//TODO: some sort of redirection to what this is supposed to lead to
	}

	//Tester Build
	if version == "t" {
		if _, err := os.Stat("osu-title-testers.png"); errors.Is(err, os.ErrNotExist) {
			ctx.Data(http.StatusOK, "image/png", []byte{})
		} else {
			titleImage, error := os.ReadFile("osu-title-testers.png")
			titleImageHash := md5.Sum(titleImage)

			if current == hex.EncodeToString(titleImageHash[:]) {
				ctx.Data(http.StatusOK, "image/png", []byte{})
				return
			}

			if error != nil {
				ctx.Data(http.StatusOK, "image/png", []byte{})
				return
			}

			ctx.Data(http.StatusOK, "image/png", titleImage)
		}
	} else {
		if _, err := os.Stat("osu-title.png"); errors.Is(err, os.ErrNotExist) {
			ctx.Data(http.StatusOK, "image/png", []byte{})
		} else {
			titleImage, error := os.ReadFile("osu-title.png")
			titleImageHash := md5.Sum(titleImage)

			if current == hex.EncodeToString(titleImageHash[:]) {
				ctx.Data(http.StatusOK, "image/png", []byte{})
				return
			}

			if error != nil {
				ctx.Data(http.StatusOK, "image/png", []byte{})
				return
			}

			ctx.Data(http.StatusOK, "image/png", titleImage)
		}
	}
}
