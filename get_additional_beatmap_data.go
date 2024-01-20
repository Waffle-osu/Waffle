package main

import (
	"Waffle/database"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func GetAdditionalBeatmapInfo() {
	//Setup Logger
	filename := fmt.Sprintf("logs/%d-log-beatmap_import.txt", time.Now().Unix())

	file, fileErr := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if fileErr != nil {
		panic(fileErr)
	}

	multiWriter := io.MultiWriter(file, os.Stdout)

	logger := log.New(multiWriter, "Beatmap Importer: ", log.LstdFlags|log.Lshortfile)

	apiKeyFile, apiKeyFileErr := ioutil.ReadFile(".api_key")

	if apiKeyFileErr != nil {
		logger.Printf("Failed to read API Key file, make sure you have a .api_key file with just your API key inside")
		return
	}

	apiKey := string(apiKeyFile)

	allBeatmapSetQuery, allBeatmapSetQueryErr := database.Database.Query("SELECT beatmapset_id FROM beatmapsets WHERE beatmapset_id < 10000000")

	if allBeatmapSetQueryErr != nil {
		logger.Fatalf("Failed to query all beatmaps!")
	}

	errors := []error{}

	setIds := []int32{}

	for allBeatmapSetQuery.Next() {
		setId := int32(0)

		allBeatmapSetQuery.Scan(&setId)

		setIds = append(setIds, setId)
	}

	allBeatmapSetQuery.Close()

	for _, setId := range setIds {
		url := fmt.Sprintf("https://osu.ppy.sh/api/get_beatmaps?k=%s&s=%d", apiKey, setId)
		response, getErr := http.Get(url)

		if getErr != nil {
			logger.Printf("Failed to ping API on Set ID %d\n", setId)
			errors = append(errors, getErr)
			continue
		}

		body, readErr := ioutil.ReadAll(response.Body)

		if readErr != nil {
			logger.Printf("Failed to read API response on Set ID %d\n", setId)
			errors = append(errors, readErr)
			continue
		}

		var beatmapInfos OsuApiBeatmapsetResponse

		jsonParseErr := json.Unmarshal(body, &beatmapInfos)

		if jsonParseErr != nil {
			logger.Printf("Failed to Parse JSON response on Set ID %d\n", setId)
			errors = append(errors, jsonParseErr)
			continue
		}

		if len(beatmapInfos) == 0 {
			logger.Printf("Nothing to process on %d\n", setId)
			continue
		}

		genreId, genreIdParseErr := strconv.ParseInt(beatmapInfos[0].GenreID, 10, 64)
		languageId, languageIdParseErr := strconv.ParseInt(beatmapInfos[0].LanguageID, 10, 64)

		if genreIdParseErr != nil || languageIdParseErr != nil {
			logger.Printf("Failed to parse genre/language id on Set ID %d\n", setId)
			errors = append(errors, genreIdParseErr)
			errors = append(errors, languageIdParseErr)
			continue
		}

		updateQuery := "UPDATE beatmapsets SET genre_id = ?, language_id = ?, beatmap_pack = ? WHERE beatmapset_id = ?"
		_, updateErr := database.Database.Exec(updateQuery, genreId, languageId, beatmapInfos[0].Packs, setId)

		if updateErr != nil {
			logger.Printf("Failed to update on Set ID %d\n", setId)
			errors = append(errors, updateErr)
			continue
		}

		insertPlaycountsQuery := "INSERT INTO osu_bancho_beatmap_playcounts (beatmap_id, beatmapset_id, passcount, playcount, mode) VALUES (?, ?, ?, ?, ?)"

		for _, element := range beatmapInfos {
			beatmapId, beatmapIdParseErr := strconv.ParseInt(element.BeatmapID, 10, 64)

			if beatmapIdParseErr != nil {
				logger.Printf("Failed to insert on Set ID %d\n", setId)
				errors = append(errors, beatmapIdParseErr)
				goto outside
			}

			playcount, playcountParseErr := strconv.ParseInt(element.Playcount, 10, 64)

			if playcountParseErr != nil {
				logger.Printf("Failed to insert on Set ID %d\n", setId)
				errors = append(errors, playcountParseErr)
				goto outside
			}

			passcount, passcountParseErr := strconv.ParseInt(element.Passcount, 10, 64)

			if passcountParseErr != nil {
				logger.Printf("Failed to insert on Set ID %d\n", setId)
				errors = append(errors, passcountParseErr)
				goto outside
			}

			mode, modeParseErr := strconv.ParseInt(element.Mode, 10, 64)

			if modeParseErr != nil {
				logger.Printf("Failed to insert on Set ID %d\n", setId)
				errors = append(errors, modeParseErr)
				goto outside
			}

			_, insertErr := database.Database.Exec(insertPlaycountsQuery, beatmapId, setId, passcount, playcount, mode)

			if insertErr != nil {
				logger.Printf("Failed to insert on Set ID %d\n", setId)
				errors = append(errors, insertErr)
				goto outside
			}
		}

	outside:

		fmt.Printf("Set ID %d successfully updated.\n", setId)

		//time.Sleep(150 * time.Millisecond)
	}

	logger.Printf("Additional data get complete.")

	if len(errors) != 0 {
		logger.Printf("Errors:")

		for index, element := range errors {
			fmt.Printf("%d - %s", index, element.Error())
		}
	}
}
