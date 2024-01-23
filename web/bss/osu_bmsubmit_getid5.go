package bss

import (
	"Waffle/database"
	"Waffle/helpers"
	"Waffle/utils"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"
	"unicode"

	"github.com/Waffle-osu/osu-parser/osu_parser"
	"github.com/gin-gonic/gin"
)

func isASCII(s string) bool {
	for _, c := range s {
		if c > unicode.MaxASCII {
			return false
		}
	}
	return true
}

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

	uploadRequest := GetUploadRequest(userId)

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

		_, exists, approved, setId, queryErrorOccured := CheckBeatmapStatus(uploadRequest.BeatmapsetId, userData, &uploadRequest.Metadata)

		if queryErrorOccured {
			ctx.String(500, "Internal Queries failed!")

			return true
		}

		osuTicketFormat := fmt.Sprintf("%d-%s-%s-%s", time.Now().Unix(), username, osuFormFile.Filename, userData.Password)
		osuTicketBytes := sha256.Sum256([]byte(osuTicketFormat))
		osuTicketHashed := osuTicketBytes[:]
		osuTicket := hex.EncodeToString(osuTicketHashed)

		filename := osuFormFile.Filename

		if !isASCII(parsedOsu.Metadata.Title) || !isASCII(parsedOsu.Metadata.Artist) {
			filename = fmt.Sprintf("%s - %s (%s) [%s].osu", parsedOsu.Metadata.Artist, parsedOsu.Metadata.Title, parsedOsu.Metadata.Creator, parsedOsu.Metadata.Version)
		}

		uploadTicket := UploadTicket{
			Ticket:    osuTicket,
			Filename:  filename,
			Size:      osuFormFile.Size,
			ParsedOsu: parsedOsu,
			FileData:  readOutOsuFile,
		}

		uploadRequest.UploadTickets[osuTicket] = uploadTicket

		helpers.Logger.Printf("BSS: UploadTicket generated for: %s", osuFormFile.Filename)

		oszFileName := fmt.Sprintf("%d %s - %s (%s).osz", uploadRequest.BeatmapsetId, uploadRequest.Metadata.Artist, uploadRequest.Metadata.Title, uploadRequest.Metadata.Creator)

		returnMessage := "new\n"

		subject := ""
		message := ""

		if exists {
			uploadRequest.IsUpdate = true
			uploadRequest.BeatmapsetId = setId

			getPostInfoSql := "SELECT subject, message FROM osu_beatmap_posts WHERE beatmapset_id = ?"
			getPostInfoQuery, getPostInfoQueryErr := database.Database.Query(getPostInfoSql, setId)

			if getPostInfoQueryErr != nil {
				ctx.String(500, "internal queries failed")

				return true
			}

			if getPostInfoQuery.Next() {
				postInfoScanErr := getPostInfoQuery.Scan(&subject, &message)
				getPostInfoQuery.Close()

				if postInfoScanErr != nil {
					ctx.String(500, "internal queries failed")

					return true
				}
			}

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
		returnMessage += fmt.Sprintf("%d\n", uploadRequest.BeatmapsetId) //Thread ID
		returnMessage += fmt.Sprintf("%s\n", formattedApproved)
		returnMessage += fmt.Sprintf("%s\n", subject) //Subject
		returnMessage += fmt.Sprintf("%s\n", message) //Message

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
			uploadRequest.IsUpdate = true

			//Get currently stored versions
			versionsGetSql := `SELECT version FROM beatmaps WHERE beatmapset_id = ?`
			versionsGetQuery, versionsGetErr := database.Database.Query(versionsGetSql, uploadRequest.BeatmapsetId)

			if versionsGetErr != nil {
				ctx.String(500, "Internal Queries failed.")

				return true
			}

			//Run diff to figure out which diffs got deleted/renamed
			currentVersions := map[string]string{}
			uploadVersions := map[string]string{}
			addedVersions := map[string]string{}
			removedVersions := map[string]string{}

			//Current Versions
			for versionsGetQuery.Next() {
				versionName := sql.NullString{}
				scanErr := versionsGetQuery.Scan(&versionName)

				if scanErr != nil {
					versionsGetQuery.Close()
					return true
				}

				if versionName.Valid {
					currentVersions[versionName.String] = versionName.String
				}
			}

			//All Currently Uploaded Versions
			for _, ticket := range uploadRequest.UploadTickets {
				uploadVersions[ticket.ParsedOsu.Metadata.Version] = ticket.ParsedOsu.Metadata.Version
			}

			//All Removed Versions
			for version := range currentVersions {
				_, exists := uploadVersions[version]

				if !exists {
					removedVersions[version] = version
				}
			}

			//All Added Versions
			for version := range uploadVersions {
				old, exists := currentVersions[version]

				if !exists {
					addedVersions[version] = old
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
			for version := range removedVersions {
				versionDeleteSql := "DELETE FROM beatmaps WHERE beatmapset_id = ? AND version = ?"
				_, versionDeleteQueryErr := database.Database.Exec(versionDeleteSql, uploadRequest.BeatmapsetId, version)

				if versionDeleteQueryErr != nil {
					ctx.String(500, "Internal queries failed.")

					return true
				}
			}

			//Update existing diffs
			for _, ticket := range uploadRequest.UploadTickets {
				diff := ticket.ParsedOsu.Difficulty

				updateDiffSql := `
					UPDATE 
						beatmaps 
					SET 
						beatmap_md5               = ?, 
						total_length              = ?, 
						drain_time                = ?, 
						count_objects             = ?, 
						count_normal              = ?,
						count_slider              = ?, 
						count_spinner             = ?, 
						diff_hp                   = ?, 
						diff_cs                   = ?, 
						diff_od                   = ?, 
						playmode                  = ?, 
						status_valid_from_version = ?,
						last_update               = CURRENT_TIMESTAMP() 
					WHERE 
						beatmapset_id = ? AND 
						version = ?
				`

				minVersion := utils.VersionOsuFile(ticket.ParsedOsu)

				_, updateDiffErr :=
					database.Database.Exec(
						updateDiffSql,
						ticket.ParsedOsu.Md5Hash,
						ticket.ParsedOsu.Length,
						ticket.ParsedOsu.DrainLength,
						len(ticket.ParsedOsu.HitObjects.List),
						ticket.ParsedOsu.HitObjects.CountNormal,
						ticket.ParsedOsu.HitObjects.CountSlider,
						ticket.ParsedOsu.HitObjects.CountSpinner,
						diff.HPDrainRate,
						diff.CircleSize,
						diff.OverallDifficulty,
						ticket.ParsedOsu.General.Mode,
						minVersion,
						uploadRequest.BeatmapsetId,
						ticket.ParsedOsu.Metadata.Version,
					)

				if updateDiffErr != nil {
					ctx.String(500, "Internal Queries Failed")

					return true
				}
			}

			//Insert the new diffs
			toInsert := []UploadTicket{}

			for version := range addedVersions {
				for _, ticket := range uploadRequest.UploadTickets {
					if ticket.ParsedOsu.Metadata.Version == version {
						toInsert = append(toInsert, ticket)
					}
				}
			}

			for _, ticket := range toInsert {
				currentOsu := ticket.ParsedOsu

				insertBeatmapErr := InsertIntoBeatmaps(currentOsu, uploadRequest.BeatmapsetId, -userId, ticket.Filename)

				if insertBeatmapErr != nil {
					ctx.String(500, "Internal queries failed")

					return true
				}
			}

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

				insertBeatmapErr := InsertIntoBeatmaps(currentOsu, uploadRequest.BeatmapsetId, -userId, ticket.Filename)

				if insertBeatmapErr != nil {
					ctx.String(500, "Internal queries failed")

					return true
				}
			}
		}

		//Copy over all the .osu files
		for _, ticket := range uploadRequest.UploadTickets {
			filename := fmt.Sprintf("osus/%s", ticket.Filename)

			if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
				os.WriteFile(filename, ticket.FileData, 0755)
			} else {
				file, err := os.Create(filename)

				if err != nil {
					ctx.String(500, "Internal queries failed")

					return true
				}

				file.Write(ticket.FileData)
				file.Close()
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

		beatmapsetId, newSetIdErr := RegisterRequest(userId, &newUploadRequest)

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
