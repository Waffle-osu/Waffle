package bss

import (
	"Waffle/database"
	"Waffle/helpers"
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
	//setId := ctx.Query("s")
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

	if !success {
		ctx.String(401, "Authentication failed!")

		return
	}

	querySuccess, userData := database.UserFromDatabaseById(uint64(userId))

	if querySuccess != 0 {
		ctx.String(401, "Authentication failed!")

		return
	}

	uploadRequest := GetUploadRequest(int64(userId))

	parsedOsu, err := osu_parser.ParseText(string(readOutOsuFile))

	if err != nil {
		ctx.String(400, "Could not read out file")

		return
	}

	if osuFileErr != nil {
		ctx.String(400, "Could not read out file")

		return
	}

	//returns true if some error occured in the function and we need to return
	commonPush := func() bool {
		if uploadRequest == nil {
			ctx.String(400, "Action is invalid in this context.")

			return true
		}

		_, exists, approved, _, queryErrorOccured := CheckBeatmapStatus(uploadRequest.BeatmapsetId, userData, &uploadRequest.Metadata)

		if queryErrorOccured != false {
			ctx.String(500, "Internal Queries failed!")

			return true
		}

		osuTicket := fmt.Sprintf("%d-%s-%s-%s", time.Now().Unix(), username, osuFormFile.Filename, userData.Password)
		osuTicketBytes := sha256.Sum256([]byte(osuTicket))
		osuTicketHashed := string(osuTicketBytes[:])

		uploadTicket := UploadTicket{
			Ticket:   osuTicketHashed,
			Filename: osuFormFile.Filename,
			Size:     osuFormFile.Size,
		}

		uploadRequest.UploadTickets = append(uploadRequest.UploadTickets, uploadTicket)

		helpers.Logger.Printf("BSS: UploadTicket generated for: %s", osuFormFile.Filename)

		oszFileName := fmt.Sprintf("%d %s - %s (%s).osz", uploadRequest.BeatmapsetId, uploadRequest.Metadata.Artist, uploadRequest.Metadata.Title, uploadRequest.Metadata.Creator)

		returnMessage := "new\n"

		if exists {
			returnMessage = "old\n"
		}

		formattedApproved := "0"

		if approved {
			formattedApproved = "1"
		}

		returnMessage += fmt.Sprintf("%d\n", uploadRequest.BeatmapsetId)
		returnMessage += fmt.Sprintf("%s\n", uploadRequest.OszTicket)
		returnMessage += fmt.Sprintf("%s\n", uploadTicket.Ticket)
		returnMessage += fmt.Sprintf("%s\n", oszFileName)
		returnMessage += fmt.Sprintf("%d\n", 0) //Thread ID
		returnMessage += fmt.Sprintf("%s\n", formattedApproved)
		returnMessage += fmt.Sprintf("%s\n", "") //Subject
		returnMessage += fmt.Sprintf("%s", "")   //Message

		ctx.String(200, returnMessage)

		return false
	}

	//same as commonPush, returns true if whole request needs to return
	commonSubmissionDone := func() bool {
		if uploadRequest == nil {
			ctx.String(400, "Action is invalid in this context.")

			return true
		}

		_, exists, _, _, queryErrorOccured := CheckBeatmapStatus(uploadRequest.BeatmapsetId, userData, &uploadRequest.Metadata)

		if queryErrorOccured != false {
			ctx.String(500, "Internal queries failed!")

			return true
		}

		if exists {
			uploadRequest.IsUpdate = exists
		} else {
			//Create Beatmapset
		}

		return false
	}

	switch action {
	//Push file
	case "0":
		if commonPush() {
			return
		}

	//Initial submission, first push
	case "1":
		if uploadRequest != nil {
			ctx.String(400, "Action is invalid in this context.")

			return
		}

		oszTicket := fmt.Sprintf("%d-%s-%s-%s-oszTicket", time.Now().Unix(), username, osuFormFile.Filename, userData.Password)
		oszTicketBytes := sha256.Sum256([]byte(oszTicket))
		oszTicketHashed := string(oszTicketBytes[:])

		newUploadRequest := UploadRequest{
			UploadTickets: []UploadTicket{},
			HasVideo:      hasVideo == "1",
			HasStoryboard: hasStoryboard == "1",
			OszTicket:     oszTicketHashed,
			Metadata:      parsedOsu.Metadata,
		}

		beatmapsetId, newSetIdErr := RegisterRequest(int64(userId), &newUploadRequest)

		uploadRequest = &newUploadRequest

		if newSetIdErr != nil {
			ctx.String(500, "Internal queries failed!")

			return
		}

		newUploadRequest.BeatmapsetId = beatmapsetId

		helpers.Logger.Printf("BSS: Created new UploadRequest for %s", username)

		if commonPush() {
			return
		}
	//Push last, ends submission
	case "2":
		if commonPush() {
			return
		}

		if commonSubmissionDone() {
			return
		}
	//Push single and end
	case "3":
		if commonPush() {
			return
		}

		if commonSubmissionDone() {
			return
		}
	}
}
