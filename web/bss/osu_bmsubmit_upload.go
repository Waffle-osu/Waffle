package bss

import (
	"Waffle/database"
	"Waffle/helpers"
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func HandleUpload(ctx *gin.Context) {
	username := ctx.Query("u")
	password := ctx.Query("p")
	ticket := ctx.Query("c")
	first := ctx.Query("r")
	oszFilename := ctx.Query("of")
	oszTicket := ctx.Query("oc")
	setId := ctx.Query("s")

	file, fileErr := ctx.FormFile("osu")

	if fileErr != nil {
		ctx.String(400, "Failed to read out file")
	}

	userId, success := database.AuthenticateUser(username, password)

	if !success {
		ctx.String(401, "Authentication failed!")

		return
	}

	uploadRequest := GetUploadRequest(int64(userId))

	//full .osz upload
	if ticket == uploadRequest.OszTicket {
		os.Remove(fmt.Sprintf("oszs/%d.osz", uploadRequest.BeatmapsetId))

		openFile, openFileErr := file.Open()

		if openFileErr != nil {
			ctx.String(500, "Failed to open osz")

			return
		}

		newOsz, createErr := os.Create(fmt.Sprintf("oszs/%d.osz", uploadRequest.BeatmapsetId))

		if createErr != nil {
			ctx.String(500, "Failed to create new osz")

			return
		}

		readOutBuffer := make([]byte, file.Size)
		openFile.Read(readOutBuffer)
		newOsz.Write(readOutBuffer)

		openFile.Close()
		newOsz.Close()

		DeleteUploadRequest(int64(userId))
	} else {
		//Extract existing osz
		if first == "1" {
			archive, archiveOpenErr := zip.OpenReader(fmt.Sprintf("oszs/%d.osz", uploadRequest.BeatmapsetId))

			if archiveOpenErr != nil {
				ctx.String(500, "Failed to read existing osz file")

				return
			}

			dirName := fmt.Sprintf("bss_temp/%d", userId)
			_, err := os.Stat(dirName)

			if err == nil {
				ctx.String(500, "BSS Upload failed.")

				return
			}

			dirCreateErr := os.Mkdir(dirName, os.ModePerm)

			if dirCreateErr == nil {
				ctx.String(500, "BSS Upload failed.")

				return
			}

			for _, f := range archive.File {
				filePath := filepath.Join(dirName, f.Name)

				//They all get reuploaded, we don't need them twice
				if strings.HasSuffix(filePath, ".osu") {
					continue
				}

				if f.FileInfo().IsDir() {
					os.MkdirAll(filePath, os.ModePerm)

					continue
				}

				if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
					panic(err)
				}

				dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
				if err != nil {
					panic(err)
				}

				fileInArchive, err := f.Open()
				if err != nil {
					panic(err)
				}

				if _, err := io.Copy(dstFile, fileInArchive); err != nil {
					panic(err)
				}

				dstFile.Close()
				fileInArchive.Close()
			}
		}

		newFile, newFileErr := os.Create(fmt.Sprintf("bss_temp/%d/%s", userId, file.Filename))

		if newFileErr != nil {
			ctx.String(500, "Failed to finish upload")

			return
		}

		openFile, openFileErr := file.Open()

		if openFileErr != nil {
			ctx.String(500, "Failed to open osu")

			return
		}

		readOutBuffer := make([]byte, file.Size)
		openFile.Read(readOutBuffer)
		newFile.Write(readOutBuffer)

		newFile.Close()
		openFile.Close()

		_, exists := uploadRequest.UploadTickets[ticket]

		if exists {
			delete(uploadRequest.UploadTickets, ticket)
		} else {
			ctx.String(400, "Invalid ticket")

			return
		}
	}

	helpers.Logger.Printf("-- Got upload:\n")
	helpers.Logger.Printf("Username: %s\n", username)
	helpers.Logger.Printf("Ticket: %s\n", ticket)
	helpers.Logger.Printf("NoIdea: %s\n", first)
	helpers.Logger.Printf("OszFilename: %s\n", oszFilename)
	helpers.Logger.Printf("OszTicket: %s\n", oszTicket)
	helpers.Logger.Printf("SetId: %s\n", setId)
	helpers.Logger.Printf("Filename: %s\n", file.Filename)
	helpers.Logger.Printf("File size: %d\n", file.Size)

	ctx.String(200, "ok")
}
