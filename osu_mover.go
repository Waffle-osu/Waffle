package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func MoveOsuFiles(osuDir string) {
	//Setup Logger
	filename := fmt.Sprintf("logs/%d-log-osu_move.txt", time.Now().Unix())

	file, fileErr := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if fileErr != nil {
		panic(fileErr)
	}

	multiWriter := io.MultiWriter(file, os.Stdout)

	logger := log.New(multiWriter, ".osz Renamer: ", log.LstdFlags|log.Lshortfile)

	songsDirectory, songsDirectoryReadErr := ioutil.ReadDir(osuDir)

	if songsDirectoryReadErr != nil {
		logger.Fatalf("Failed to read osz directory!")
	}

	for _, directory := range songsDirectory {
		if !directory.IsDir() {
			continue
		}

		folderFiles, folderFilesErr := ioutil.ReadDir(osuDir + "/" + directory.Name())

		if folderFilesErr != nil {
			logger.Fatalf("Failed to read directory %s", directory.Name())
		}

		for _, file := range folderFiles {
			if file.IsDir() {
				continue
			}

			if !strings.HasSuffix(file.Name(), ".osu") {
				continue
			}

			currentFilename := osuDir + "/" + directory.Name() + "/" + file.Name()

			osuFile, osuFileOpenErr := os.Open(currentFilename)

			if osuFileOpenErr != nil {
				logger.Fatalf("Failed to open file %s", currentFilename)
			}

			defer osuFile.Close()

			out, osuFileCreateErr := os.Create("osus/" + file.Name())

			if osuFileCreateErr != nil {
				logger.Fatalf("Failed to create file %s", "osus/"+file.Name())
			}

			defer out.Close()

			_, osuFileCopyErr := io.Copy(out, osuFile)

			if osuFileCopyErr != nil {
				logger.Fatalf("Failed to copy file %s", currentFilename)
			}

			out.Close()
		}
	}
}
