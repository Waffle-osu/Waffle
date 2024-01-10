package bss

import (
	"Waffle/database"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/Waffle-osu/osu-parser/osu_parser"
	"github.com/gin-gonic/gin"
)

func HandleGetId5(ctx *gin.Context) {
	username := ctx.Query("u")
	password := ctx.Query("p")
	action := ctx.Query("r")
	setId := ctx.Query("s")
	hasVideo := ctx.Query("v")
	hasStoryboard := ctx.Query("sb")

	osuFormFile, formFileErr := ctx.FormFile("osu")

	if formFileErr != nil {
		ctx.String(400, "Could not read out file")

		return
	}

	osuFile, osuFileErr := osuFormFile.Open()
	readOutOsuFile := make([]byte, osuFormFile.Size)

	_, readErr := osuFile.Read(readOutOsuFile)

	if readErr != nil {
		ctx.String(400, "Could not read out file")

		return
	}

	userId, success := database.AuthenticateUser(username, password)

	uploadRequest := GetUploadRequest(userId)

	parsedOsu, err := osu_parser.ParseText(string(readOutOsuFile))

	if err != nil {
		ctx.String(400, "Could not read out file")

		return
	}

	if !success {
		ctx.String(401, "Authentication failed!")

		return
	}

	if osuFileErr != nil {
		ctx.String(400, "Could not read out file")

		return
	}

	switch action {
	//Push file
	case "0":
	push:
		fmt.Sprintf("Hi")

	//Initial submission, first push
	case "1":
		if uploadRequest != nil {
			ctx.String("Action is invalid in this context.")

			return
		}

		oszTicket := fmt.Sprintf("%d-%s-%s", time.Now().Unix(), username, osuFormFile.Filename)
		oszTicketBytes := sha256.Sum256([]byte(oszTicket))
		oszTicketHashed := string(oszTicketBytes[:])

		newUploadRequest := UploadRequest{
			UploadTickets: []UploadTicket{},
			HasVideo:      hasVideo == "1",
			HasStoryboard: hasStoryboard == "1",
			OszTicket:     oszTicketHashed,
			Metadata:      parsedOsu.Metadata,
		}

		goto push
	//Push last, ends submission
	case "2":
		goto push

	}
}
