package bss

import (
	"Waffle/bancho/client_manager"
	"Waffle/database"
	"Waffle/helpers"
	"Waffle/utils/zip_utils"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func HandleUpload(ctx *gin.Context) {
	username := ctx.Query("u")
	password := ctx.Query("p")
	ticket := ctx.Query("c")
	first := ctx.Query("r")
	// oszFilename := ctx.Query("of")
	// oszTicket := ctx.Query("oc")
	// setId := ctx.Query("s")

	file, fileErr := ctx.FormFile("osu")

	if fileErr != nil {
		ctx.String(400, "Failed to read out file")
	}

	userId, success := database.AuthenticateUser(username, password)

	if !success {
		ctx.String(401, "Authentication failed!")

		helpers.Logger.Printf("BSS Request for %s failed to authenticate.", username)

		return
	}

	uploadRequest := GetUploadRequest(userId)

	if uploadRequest == nil {
		ctx.String(400, "No upload request!")

		helpers.Logger.Printf("Upload request for UserID %d not found.", userId)

		return
	}

	//full .osz upload
	if ticket == uploadRequest.OszTicket {
		openFile, openFileErr := file.Open()

		if openFileErr != nil {
			ctx.String(500, "Failed to open osz")

			DeleteUploadRequest(userId)

			helpers.Logger.Printf("BSS:U %d failed to open osz", userId)

			return
		}

		newOszFilename := fmt.Sprintf("bss_temp/%d.osz", uploadRequest.BeatmapsetId)
		newOsz, createErr := os.Create(newOszFilename)

		if createErr != nil {
			ctx.String(500, "Failed to create new osz")

			helpers.Logger.Printf("BSS:U %d failed to create osz", userId)

			DeleteUploadRequest(userId)

			return
		}

		readOutBuffer := make([]byte, file.Size)
		openFile.Read(readOutBuffer)
		newOsz.Write(readOutBuffer)

		openFile.Close()
		newOsz.Close()

		tempOszDir := fmt.Sprintf("bss_temp/oszs/%d", uploadRequest.BeatmapsetId)
		os.MkdirAll(tempOszDir, 0777)
		unzipErr := zip_utils.UnzipFile(newOszFilename, tempOszDir, false)

		if unzipErr != nil {
			if errors.Is(unzipErr, errors.ErrUnsupported) {
				foundClient := client_manager.ClientManager.GetClientById(userId)

				if foundClient != nil {
					foundClient.BanchoAnnounce("Your upload failed due to a invalid filename, this could be because of Unicode characters in your Metadata. Make sure a Non-Unicode metadata is present, or if using an older client, only the regular alphabet is used.")
				}

				helpers.Logger.Printf("BSS:U %d failed due to unicode", userId)
			} else {
				helpers.Logger.Printf("BSS:U %d failed to unzip", userId)
			}

			ctx.String(500, "Failed to create new osz")

			DeleteUploadRequest(userId)

			return
		}

		tempDir, tempDirErr := os.ReadDir(tempOszDir)

		if tempDirErr != nil {
			ctx.String(500, "Failed to check osz")

			helpers.Logger.Printf("BSS:U %d failed to list files for osz", userId)

			return
		}

		//Check every .osu file and check it againt its ticket
		for _, entry := range tempDir {
			foundTicket := false

			if entry.IsDir() {
				continue
			}

			filename := entry.Name()

			if !strings.HasSuffix(filename, ".osu") {
				continue
			}

			for _, ticket := range uploadRequest.UploadTickets {
				if ticket.Filename == filename {
					foundTicket = true

					break
				}
			}

			if !foundTicket {
				ctx.String(400, "File not in ticket")

				helpers.Logger.Printf("BSS:U %d file was not in ticket", userId)

				DeleteUploadRequest(userId)

				return
			}
		}

		//Save osz
		os.Remove(fmt.Sprintf("oszs/%d.osz", uploadRequest.BeatmapsetId))

		savedOszErr := zip_utils.ZipDirectory(fmt.Sprintf("oszs/%d.osz", uploadRequest.BeatmapsetId), tempOszDir)

		if savedOszErr != nil {
			ctx.String(500, "Failed to save osz")

			helpers.Logger.Printf("BSS:U %d failed to save osz", userId)

			DeleteUploadRequest(userId)

			return
		}

		//if everything's fine, take the first ticket sent
		var firstTicket UploadTicket

		for _, value := range uploadRequest.UploadTickets {
			firstTicket = value
			break
		}

		//Generate Thumbnail and Mp3
		GenerateThumbnail(firstTicket, uploadRequest.BeatmapsetId)
		CreateMp3Preview(firstTicket.ParsedOsu.General.AudioFilename, firstTicket.ParsedOsu.General.PreviewTime, uploadRequest.BeatmapsetId)

		DeleteUploadRequest(userId)
	} else {
		_, exists := uploadRequest.UploadTickets[ticket]

		if !exists {
			ctx.String(400, "Invalid ticket")

			helpers.Logger.Printf("BSS:U %d no ticket", userId)

			return
		}

		//Extract existing osz
		if first == "1" {
			existingOsz := fmt.Sprintf("oszs/%d.osz", uploadRequest.BeatmapsetId)
			destFolder := fmt.Sprintf("bss_temp/%d", userId)

			unzipErr := zip_utils.UnzipFile(existingOsz, destFolder, true)

			if unzipErr != nil {
				ctx.String(500, "Failed to finish upload")

				helpers.Logger.Printf("BSS:U %d failed to unzip osz", userId)

				DeleteUploadRequest(userId)

				return
			}
		}

		var uploadTicket UploadTicket

		for _, listTicket := range uploadRequest.UploadTickets {
			if listTicket.Ticket == ticket {
				uploadTicket = listTicket
			}
		}

		osuFilename := uploadTicket.Filename

		if !isASCII(uploadRequest.Metadata.Artist) || !isASCII(uploadRequest.Metadata.Title) {
			osuFilename = fmt.Sprintf("%s - %s (%s) [%s].osu", uploadTicket.ParsedOsu.Metadata.Artist, uploadTicket.ParsedOsu.Metadata.Title, uploadTicket.ParsedOsu.Metadata.Creator, uploadTicket.ParsedOsu.Metadata.Version)
		}

		newFile, newFileErr := os.Create(fmt.Sprintf("bss_temp/%d/%s", userId, osuFilename))

		if newFileErr != nil {
			ctx.String(500, "Failed to finish upload")

			helpers.Logger.Printf("BSS:U %d failed to create osu file", userId)

			DeleteUploadRequest(userId)

			return
		}

		openFile, openFileErr := file.Open()

		if openFileErr != nil {
			ctx.String(500, "Failed to open osu")

			helpers.Logger.Printf("BSS:U %d failed to open osu file", userId)

			DeleteUploadRequest(userId)

			return
		}

		readOutBuffer := make([]byte, file.Size)
		openFile.Read(readOutBuffer)
		newFile.Write(readOutBuffer)

		newFile.Close()
		openFile.Close()

		delete(uploadRequest.UploadTickets, ticket)

		//All tickets gone, create osz
		if len(uploadRequest.UploadTickets) == 0 {
			oszFilename := fmt.Sprintf("bss_temp/%d.osz", uploadRequest.BeatmapsetId)

			oszCreateErr := zip_utils.ZipDirectory(oszFilename, fmt.Sprintf("bss_temp/%d", userId))

			if oszCreateErr != nil {
				ctx.String(500, "Failed to finish upload")

				helpers.Logger.Printf("BSS:U %d failed to create new osz after update", userId)

				DeleteUploadRequest(userId)

				return
			}

			oldPath := fmt.Sprintf("oszs/%d.osz", uploadRequest.BeatmapsetId)

			removeErr := os.Remove(oldPath)

			if removeErr != nil && !errors.Is(removeErr, os.ErrNotExist) {
				ctx.String(500, "Failed to finish upload")

				helpers.Logger.Printf("BSS:U %d failed to remove old osz", userId)

				os.Remove(oszFilename)

				return
			}

			os.Rename(oszFilename, oldPath)

			DeleteUploadRequest(userId)
		}
	}

	ctx.String(200, "ok")
}
