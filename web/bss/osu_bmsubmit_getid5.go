package bss

import (
	"Waffle/database"
	"Waffle/helpers"
	"Waffle/utils"
	"crypto/md5"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
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
	osuFile.Close()

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

		if queryErrorOccured {
			ctx.String(500, "Internal Queries failed!")

			return true
		}

		osuTicketFormat := fmt.Sprintf("%d-%s-%s-%s", time.Now().Unix(), username, osuFormFile.Filename, userData.Password)
		osuTicketBytes := sha256.Sum256([]byte(osuTicketFormat))
		osuTicketHashed := osuTicketBytes[:]
		osuTicket := hex.EncodeToString(osuTicketHashed)

		uploadTicket := UploadTicket{
			Ticket:    osuTicket,
			Filename:  osuFormFile.Filename,
			Size:      osuFormFile.Size,
			ParsedOsu: parsedOsu,
			FileData:  readOutOsuFile,
		}

		uploadRequest.UploadTickets[osuTicket] = uploadTicket

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

		if queryErrorOccured {
			ctx.String(500, "Internal queries failed!")

			return true
		}

		if exists {
			uploadRequest.IsUpdate = exists

			//Get currently stored versions
			versionsGetSql := `SELECT version FROM beatmaps WHERE beatmapset_id = ?`
			versionsGetQuery, versionsGetErr := database.Database.Query(versionsGetSql, uploadRequest.BeatmapsetId)

			if versionsGetErr != nil {
				ctx.String(500, "Internal Queries failed.")

				return true
			}

			//Run diff to figure out which diffs got deleted/renamed
			currentVersions := map[string]bool{}
			uploadVersions := map[string]bool{}
			removedVersions := map[string]bool{}
			addedVersions := map[string]bool{}

			//Current Versions
			for versionsGetQuery.Next() {
				versionName := sql.NullString{}
				scanErr := versionsGetQuery.Scan(&versionName)

				if scanErr != nil {
					versionsGetQuery.Close()
					return true
				}

				if versionName.Valid {
					currentVersions[versionName.String] = true
				}
			}

			//All Currently Uploaded Versions
			for _, ticket := range uploadRequest.UploadTickets {
				uploadVersions[ticket.ParsedOsu.Metadata.Version] = true
			}

			//All Removed Versions
			for version, _ := range uploadVersions {
				_, exists := currentVersions[version]

				if !exists {
					removedVersions[version] = true
				}
			}

			//All Added Versions
			for version, _ := range currentVersions {
				_, exists := uploadVersions[version]

				if !exists {
					addedVersions[version] = true
				}
			}

			if len(parsedOsu.TimingPoints.TimingPoints) == 0 {
				ctx.String(400, "Invalid Timing")

				return true
			}

			bpm := 60000.0 / parsedOsu.TimingPoints.TimingPoints[0].BeatLength

			//Update metadata as necessary
			metadataUpdateSql := "UPDATE beatmapsets SET artist = ?, title = ?, creator = ?, source = ?, tags = ?, has_video = ?, has_storyboard = ?, bpm = ? WHERE beatmapset_id = ?"
			_, metadataUpdateErr := database.Database.Exec(metadataUpdateSql, uploadRequest.Metadata.Artist, uploadRequest.Metadata.Title, uploadRequest.Metadata.Creator, uploadRequest.Metadata.Source, uploadRequest.Metadata.Tags, uploadRequest.HasVideo, uploadRequest.HasStoryboard, bpm, uploadRequest.BeatmapsetId)

			if metadataUpdateErr != nil {
				ctx.String(500, "Internal queries failed.")

				return true
			}

			//Delete the Deleted diffs
			//Insert the new diffs

		} else {
			//Create Beatmapset
			bpm := 60000.0 / parsedOsu.TimingPoints.TimingPoints[0].BeatLength

			insertBeatmapsetSql := "INSERT INTO beatmapsets (beatmapset_id, creator_id, artist, title, creator, source, tags, has_video, has_storyboard, bpm) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
			_, insertBeatmapsetSqlErr := database.Database.Exec(insertBeatmapsetSql, uploadRequest.BeatmapsetId, -userId, uploadRequest.Metadata.Artist, uploadRequest.Metadata.Title, uploadRequest.Metadata.Creator, uploadRequest.Metadata.Source, uploadRequest.Metadata.Tags, uploadRequest.HasVideo, uploadRequest.HasStoryboard, bpm)

			if insertBeatmapsetSqlErr != nil {
				ctx.String(500, "Internal queries Failed!")

				return true
			}

			for _, ticket := range uploadRequest.UploadTickets {
				newBeatmapId := GetNextBssBeatmapId()

				if newBeatmapId == -1 {
					ctx.String(500, "Internal queries failed!")

					return true
				}

				currentOsu := ticket.ParsedOsu

				fileHashedBytes := md5.Sum(ticket.FileData)
				fileHashedHex := hex.EncodeToString(fileHashedBytes[:])

				hitObjectCount := len(ticket.ParsedOsu.HitObjects.HitObjects)
				totalLength := currentOsu.HitObjects.HitObjects[hitObjectCount-1].Time - currentOsu.HitObjects.HitObjects[0].Time

				breakTime := int32(0)

				for _, event := range ticket.ParsedOsu.Events.Events {
					if event.EventType == osu_parser.EventTypeBreak {
						breakTime += event.BreakTimeEnd - event.BreakTimeBegin
					}
				}

				countNormal := int32(0)
				countSlider := int32(0)
				countSpinner := int32(0)

				for _, object := range currentOsu.HitObjects.HitObjects {
					if (int32(object.Type) & int32(osu_parser.HitObjectTypeCircle)) > 0 {
						countNormal++
					}

					if (int32(object.Type) & int32(osu_parser.HitObjectTypeSlider)) > 0 {
						countSlider++
					}

					if (int32(object.Type) & int32(osu_parser.HitObjectTypeSpinner)) > 0 {
						countSpinner++
					}

					if (int32(object.Type) & int32(osu_parser.HitObjectTypeHold)) > 0 {
						countSlider++
					}
				}

				drainTime := totalLength - float64(breakTime)

				minVersion := utils.VersionOsuFile(currentOsu)

				insertBeatmapSql := "INSERT INTO beatmaps (beatmap_id, beatmapset_id, creator_id, filename, beatmap_md5, version, total_length, drain_time, count_objects, count_normal, count_slider, count_spinner, diff_hp, diff_cs, diff_od, diff_stars, playmode, ranking_status, last_update, submit_date, approve_date, beatmap_source, status_valid_from_version, status_valid_to_version) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP(), CURRENT_TIMESTAMP(), '1000-01-01 00:00:00.000000', ?, ?, ?)"
				_, insertBeatmapErr := database.Database.Query(insertBeatmapSql, newBeatmapId, uploadRequest.BeatmapsetId, -userId, ticket.Filename, fileHashedHex, currentOsu.Metadata.Version, totalLength, drainTime, hitObjectCount, countNormal, countSlider, countSpinner, currentOsu.Difficulty.HPDrainRate, currentOsu.Difficulty.CircleSize, currentOsu.Difficulty.OverallDifficulty, -1, byte(currentOsu.General.Mode), 0, 1, minVersion, 99999999)

				if insertBeatmapErr != nil {
					ctx.String(500, "Internal queries failed")

					return true
				}
			}
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

		oszTicketFormat := fmt.Sprintf("%d-%s-%s-%s-oszTicket", time.Now().Unix(), username, osuFormFile.Filename, userData.Password)
		oszTicketBytes := sha256.Sum256([]byte(oszTicketFormat))
		oszTicketHashed := oszTicketBytes[:]
		oszTicket := hex.EncodeToString([]byte(oszTicketHashed))

		newUploadRequest := UploadRequest{
			UploadTickets: map[string]UploadTicket{},
			HasVideo:      hasVideo == "1",
			HasStoryboard: hasStoryboard == "1",
			OszTicket:     oszTicket,
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
