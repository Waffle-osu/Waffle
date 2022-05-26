package web

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func HandleOsuGetDirectMp3Preview(ctx *gin.Context) {
	filename := ctx.Param("filename")
	filename = strings.TrimSuffix(filename, ".mp3")

	if _, err := os.Stat("mp3_previews/" + filename); errors.Is(err, os.ErrNotExist) {
		mp3response, getErr := http.Get("https://s.ppy.sh/mp3/preview/" + filename)

		if getErr != nil {
			ctx.Data(http.StatusOK, "waffle/blob", []byte{})
			return
		}

		outputFile, outputFileErr := os.Create("mp3_previews/" + filename)

		if outputFileErr != nil {
			ctx.Data(http.StatusOK, "waffle/blob", []byte{})

			mp3response.Body.Close()

			return
		}

		_, copyErr := io.Copy(outputFile, mp3response.Body)

		if copyErr != nil {
			outputFile.Close()
			mp3response.Body.Close()

			ctx.Data(http.StatusOK, "waffle/blob", []byte{})
			return
		}

		outputFile.Close()
		mp3response.Body.Close()

		mp3, error := os.ReadFile("mp3_previews/" + filename)

		if error != nil {
			ctx.Data(http.StatusOK, "waffle/blob", []byte{})
			return
		}

		ctx.Data(http.StatusOK, "waffle/blob", mp3)
	} else {
		mp3, error := os.ReadFile("mp3_previews/" + filename)

		if error != nil {
			ctx.Data(http.StatusOK, "waffle/blob", []byte{})
			return
		}

		ctx.Data(http.StatusOK, "waffle/blob", mp3)
	}
}
