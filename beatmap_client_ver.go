package main

import (
	"Waffle/utils"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Waffle-osu/osu-parser/osu_parser"
)

type VersionedFile struct {
	Filename          string
	DeterminedVersion int64
}

type ProcessedOsuFile struct {
	OsuFile  osu_parser.OsuFile
	Filename string
}

func RunBeatmapClientVersionDetector(osuFileDir string) {
	//Setup Logger
	filename := fmt.Sprintf("logs/%d-log-beatmap_import.txt", time.Now().Unix())

	file, fileErr := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if fileErr != nil {
		panic(fileErr)
	}

	multiWriter := io.MultiWriter(file, os.Stdout)

	logger := log.New(multiWriter, "Beatmap Versioner: ", log.LstdFlags|log.Lshortfile)

	fileEntries, readDirErr := os.ReadDir(osuFileDir)

	if readDirErr != nil {
		logger.Fatalf("Failed to read directory.")
	}

	now := time.Now()

	completedOsus := make(chan ProcessedOsuFile, len(fileEntries))
	completedVersions := make(chan VersionedFile, len(fileEntries))
	waitGroupOsus := sync.WaitGroup{}
	waitGroupVersions := sync.WaitGroup{}

	waitGroupOsus.Add(len(fileEntries))
	waitGroupVersions.Add(len(fileEntries))

	for _, file := range fileEntries {
		go func(filename string) {
			osu, parseErr := osu_parser.ParseFile(fmt.Sprintf("%s/%s", osuFileDir, filename))

			if parseErr != nil {
				logger.Fatalf("Failed to parse .osu file: %s", filename)
			}

			completedOsus <- ProcessedOsuFile{
				OsuFile:  osu,
				Filename: filename,
			}

			waitGroupOsus.Done()
		}(file.Name())
	}

	waitGroupOsus.Wait()

	for len(completedOsus) != 0 {
		toProcess := <-completedOsus

		go func(osu ProcessedOsuFile) {
			determined := utils.VersionOsuFile(osu.OsuFile)

			completedVersions <- VersionedFile{
				Filename:          osu.Filename,
				DeterminedVersion: determined,
			}

			waitGroupVersions.Done()
		}(toProcess)
	}

	waitGroupVersions.Wait()

	logger.Printf("Beatmap Versioning took %d milliseconds. Beatmaps Processed: %d", time.Since(now).Milliseconds(), len(completedVersions))
}
