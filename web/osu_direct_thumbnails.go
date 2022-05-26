package web

import (
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func HandleOsuGetDirectThumbnail(ctx *gin.Context) {
	filename := ctx.Param("filename")

	if _, err := os.Stat("direct_thumbnails/" + filename); errors.Is(err, os.ErrNotExist) {
		thumbnailResponse, getErr := http.Get("https://s.ppy.sh/mt/" + filename)

		if getErr != nil {
			ctx.Data(http.StatusOK, "image/jpg", []byte{})
			return
		}

		outputFile, outputFileErr := os.Create("direct_thumbnails/" + filename)

		if outputFileErr != nil {
			ctx.Data(http.StatusOK, "image/jpg", []byte{})

			thumbnailResponse.Body.Close()

			return
		}

		_, copyErr := io.Copy(outputFile, thumbnailResponse.Body)

		if copyErr != nil {
			ctx.Data(http.StatusOK, "image/jpg", []byte{})

			outputFile.Close()
			thumbnailResponse.Body.Close()

			return
		}

		thumbnail, error := os.ReadFile("direct_thumbnails/" + filename)

		if error != nil {
			ctx.Data(http.StatusOK, "image/jpg", []byte{})
			return
		}

		ctx.Data(http.StatusOK, "image/jpg", thumbnail)
	} else {
		thumbnail, error := os.ReadFile("direct_thumbnails/" + filename)

		if error != nil {
			ctx.Data(http.StatusOK, "image/jpg", []byte{})
			return
		}

		ctx.Data(http.StatusOK, "image/jpg", thumbnail)
	}
}
