package main

import (
	"Waffle/database"
	"Waffle/helpers"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/Waffle-osu/osu-parser/osu_parser"
	"github.com/Waffle-osu/waffle-pp/difficulty"
)

type BeatmapProcessQueryResult struct {
	BeatmapId    int32
	BeatmapSetId int32
	Filename     string
	Playmode     byte
}

func InitialDiffCalcEyup() {
	//Setup Logger
	filename := fmt.Sprintf("logs/%d-log-eyup_initial_calc.txt", time.Now().Unix())

	file, fileErr := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if fileErr != nil {
		panic(fileErr)
	}

	multiWriter := io.MultiWriter(file, os.Stdout)

	logger := log.New(multiWriter, "Difficulty Calc (eyup): ", log.LstdFlags|log.Lshortfile)

	allMapsNotInDiffTableSql := `
		SELECT beatmap_id, beatmapset_id, filename, playmode FROM beatmaps WHERE beatmap_id NOT IN (
			SELECT beatmap_id FROM osu_beatmap_difficulty
		)
	`

	startTime := time.Now()

	mapsQuery, mapsQueryErr := database.Database.Query(allMapsNotInDiffTableSql)

	if mapsQueryErr != nil {
		logger.Fatalf("Could not query maps for processing... %s\n", mapsQueryErr.Error())
	}

	insertMap := func(beatmapId int32, setId int32, filename string, eyupStars float64, playmode byte) {
		_, insertErr := database.Database.Exec("INSERT INTO osu_beatmap_difficulty (beatmap_id, beatmapset_id, mode, eyup_stars) VALUES (?, ?, ?, ?)", beatmapId, setId, playmode, eyupStars)

		if insertErr != nil {
			logger.Printf("Failed to insert b: %d; s: %d; m: %d", beatmapId, setId, playmode)
		} else {
			logger.Printf("[%d %s %s]: %.2f", beatmapId, filename, helpers.FormatPlaymodes(playmode), eyupStars)
		}
	}

	for mapsQuery.Next() {
		queryResult := BeatmapProcessQueryResult{}

		scanErr := mapsQuery.Scan(&queryResult.BeatmapId, &queryResult.BeatmapSetId, &queryResult.Filename, &queryResult.Playmode)

		if scanErr != nil {
			logger.Fatalf("Failed to scan! %s\n", scanErr.Error())
		}

		parsedFile, parseErr := osu_parser.ParseFile(fmt.Sprintf("osus/%s", queryResult.Filename))

		if parseErr != nil {
			logger.Printf("Failed to parse %s\n", queryResult.Filename)
		}

		eyupStars := difficulty.CalculateEyupStars(parsedFile)

		if queryResult.Playmode == 0 {
			//Write all 4 modes, since osu! converts exist
			insertMap(queryResult.BeatmapId, queryResult.BeatmapSetId, queryResult.Filename, eyupStars, 0)
			insertMap(queryResult.BeatmapId, queryResult.BeatmapSetId, queryResult.Filename, eyupStars, 1)
			insertMap(queryResult.BeatmapId, queryResult.BeatmapSetId, queryResult.Filename, eyupStars, 2)
			insertMap(queryResult.BeatmapId, queryResult.BeatmapSetId, queryResult.Filename, eyupStars, 3)
		} else {
			//only write the mode, since for other modes converts dont exist
			insertMap(queryResult.BeatmapId, queryResult.BeatmapSetId, queryResult.Filename, eyupStars, queryResult.Playmode)
		}
	}

	logger.Printf("Took %.2f seconds for entire diffcalc", time.Since(startTime).Seconds())
}
