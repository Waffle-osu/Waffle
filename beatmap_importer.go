package main

import (
	"Waffle/database"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type OsuApiBeatmapsetResponse []struct {
	BeatmapsetID        string `json:"beatmapset_id"`
	BeatmapID           string `json:"beatmap_id"`
	Approved            string `json:"approved"`
	TotalLength         string `json:"total_length"`
	HitLength           string `json:"hit_length"`
	Version             string `json:"version"`
	FileMd5             string `json:"file_md5"`
	DiffSize            string `json:"diff_size"`
	DiffOverall         string `json:"diff_overall"`
	DiffApproach        string `json:"diff_approach"`
	DiffDrain           string `json:"diff_drain"`
	Mode                string `json:"mode"`
	CountNormal         string `json:"count_normal"`
	CountSlider         string `json:"count_slider"`
	CountSpinner        string `json:"count_spinner"`
	SubmitDate          string `json:"submit_date"`
	ApprovedDate        string `json:"approved_date"`
	LastUpdate          string `json:"last_update"`
	Artist              string `json:"artist"`
	ArtistUnicode       string `json:"artist_unicode"`
	Title               string `json:"title"`
	TitleUnicode        string `json:"title_unicode"`
	Creator             string `json:"creator"`
	CreatorID           string `json:"creator_id"`
	Bpm                 string `json:"bpm"`
	Source              string `json:"source"`
	Tags                string `json:"tags"`
	GenreID             string `json:"genre_id"`
	LanguageID          string `json:"language_id"`
	FavouriteCount      string `json:"favourite_count"`
	Rating              string `json:"rating"`
	Storyboard          string `json:"storyboard"`
	Video               string `json:"video"`
	DownloadUnavailable string `json:"download_unavailable"`
	AudioUnavailable    string `json:"audio_unavailable"`
	Playcount           string `json:"playcount"`
	Passcount           string `json:"passcount"`
	Packs               string `json:"packs"`
	MaxCombo            string `json:"max_combo"`
	DiffAim             string `json:"diff_aim"`
	DiffSpeed           string `json:"diff_speed"`
	Difficultyrating    string `json:"difficultyrating"`
}

func BeatmapImporter(songsDir string) {
	//Setup Logger
	filename := fmt.Sprintf("logs/%d-log-beatmap_import.txt", time.Now().Unix())

	file, fileErr := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if fileErr != nil {
		panic(fileErr)
	}

	multiWriter := io.MultiWriter(file, os.Stdout)

	logger := log.New(multiWriter, "Beatmap Importer: ", log.LstdFlags|log.Lshortfile)

	beatmapSetDirectories, beatmapDirectoryReadErr := ioutil.ReadDir(songsDir)
	apiKeyFile, apiKeyFileErr := ioutil.ReadFile(".api_key")

	if beatmapDirectoryReadErr != nil {
		logger.Printf("Failed to import beatmaps: %s", beatmapDirectoryReadErr.Error())
	}

	if apiKeyFileErr != nil {
		logger.Printf("Failed to read API Key file, make sure you have a .api_key file with just your API key inside")
		return
	}

	apiKey := string(apiKeyFile)

	var onlyReprocess map[string]bool = nil

	_, reprocessFileErr := os.Stat(".reprocess")

	//Reprocess file exists
	if reprocessFileErr == nil {
		reprocessFile, readFileErr := ioutil.ReadFile(".reprocess")

		if readFileErr != nil {
			logger.Fatalf("Failed to read .reprocess")
		}

		onlyReprocess = make(map[string]bool)

		for _, setId := range strings.Split(string(reprocessFile), "\n") {
			onlyReprocess[strings.ReplaceAll(setId, "\r", "")] = true
		}
	}

	for _, directory := range beatmapSetDirectories {
		errors := 0
		startTime := time.Now()

		if directory.IsDir() == false {
			continue
		}

		splitName := strings.Split(directory.Name(), " ")
		setId := splitName[0]

		if onlyReprocess != nil {
			process, exists := onlyReprocess[setId]

			if exists == false || process == false {
				continue
			}
		}

		checksumToFilename := make(map[string]string)

		folderFiles, folderFilesErr := ioutil.ReadDir(songsDir + "/" + directory.Name())

		if folderFilesErr != nil {
			logger.Printf("Failed to read files off Set ID %s", setId)
			errors++
			continue
		}

		for _, beatmapFolderFile := range folderFiles {
			if strings.HasSuffix(beatmapFolderFile.Name(), ".osu") == false {
				continue
			}

			osuFileBytes, osuReadErr := ioutil.ReadFile(songsDir + "/" + directory.Name() + "/" + beatmapFolderFile.Name())

			if osuReadErr != nil {
				logger.Printf("Failed to read osu file off Set ID %s", setId)
				errors++
				continue
			}

			osuFileHashed := md5.Sum([]byte(osuFileBytes))
			osuFileHashedString := hex.EncodeToString(osuFileHashed[:])

			checksumToFilename[osuFileHashedString] = beatmapFolderFile.Name()
		}

		url := fmt.Sprintf("https://osu.ppy.sh/api/get_beatmaps?k=%s&s=%s", apiKey, setId)
		response, getErr := http.Get(url)

		if getErr != nil {
			logger.Printf("Failed to ping API on Set ID %s", setId)
			errors++
			continue
		}

		body, readErr := ioutil.ReadAll(response.Body)

		if readErr != nil {
			logger.Printf("Failed to read API response on Set ID %s", setId)
			errors++
			continue
		}

		var beatmapInfos OsuApiBeatmapsetResponse

		jsonParseErr := json.Unmarshal(body, &beatmapInfos)

		if jsonParseErr != nil {
			logger.Printf("Failed to Parse JSON response on Set ID %s", setId)
			errors++
			continue
		}

		var currentBeatmapset *database.Beatmapset = nil
		beatmapsetBeatmaps := []database.Beatmap{}

		for _, beatmapInfo := range beatmapInfos {
			if currentBeatmapset == nil {
				currentBeatmapset = new(database.Beatmapset)

				beatmapsetId, parseErr := strconv.ParseInt(beatmapInfo.BeatmapsetID, 10, 32)
				creatorId, parseErr := strconv.ParseInt(beatmapInfo.CreatorID, 10, 32)
				hasVideo, parseErr := strconv.ParseInt(beatmapInfo.Video, 10, 32)
				hasStoryboard, parseErr := strconv.ParseInt(beatmapInfo.Storyboard, 10, 32)
				bpm, parseErr := strconv.ParseFloat(beatmapInfo.Bpm, 64)

				if parseErr != nil {
					logger.Printf("Failed to parse JSON values to their types. Set ID %s", setId)
					errors++
					continue
				}

				currentBeatmapset.BeatmapsetId = int32(beatmapsetId)
				currentBeatmapset.CreatorId = creatorId
				currentBeatmapset.HasVideo = int8(hasVideo)
				currentBeatmapset.HasStoryboard = int8(hasStoryboard)
				currentBeatmapset.Bpm = float32(bpm)

				currentBeatmapset.Artist = beatmapInfo.Artist
				currentBeatmapset.Title = beatmapInfo.Title
				currentBeatmapset.Creator = beatmapInfo.Creator
				currentBeatmapset.Source = beatmapInfo.Source
				currentBeatmapset.Tags = beatmapInfo.Tags
			}

			beatmapId, parseErr := strconv.ParseInt(beatmapInfo.BeatmapID, 10, 32)
			totalLength, parseErr := strconv.ParseInt(beatmapInfo.TotalLength, 10, 32)
			drainTime, parseErr := strconv.ParseInt(beatmapInfo.HitLength, 10, 32)
			countNormal, parseErr := strconv.ParseInt(beatmapInfo.CountNormal, 10, 32)
			countSliders, parseErr := strconv.ParseInt(beatmapInfo.CountSlider, 10, 32)
			countSpinners, parseErr := strconv.ParseInt(beatmapInfo.CountSpinner, 10, 32)
			diffHp, parseErr := strconv.ParseInt(beatmapInfo.DiffDrain, 10, 32)
			diffCs, parseErr := strconv.ParseInt(beatmapInfo.DiffSize, 10, 32)
			diffOd, parseErr := strconv.ParseInt(beatmapInfo.DiffOverall, 10, 32)
			playmode, parseErr := strconv.ParseInt(beatmapInfo.Mode, 10, 32)
			rankingStatus, parseErr := strconv.ParseInt(beatmapInfo.Approved, 10, 32)

			if parseErr != nil {
				logger.Printf("Failed to parse JSON values to their types. Set ID %s", setId)
				errors++
				continue
			}

			countObjects := countNormal + countSpinners + countSliders

			foundFilename, filenameExists := checksumToFilename[beatmapInfo.FileMd5]

			if filenameExists == false {
				logger.Printf("Failed to find matching Filename for .osu file. Set ID %s", setId)
				errors++
				continue
			}

			beatmapsetBeatmaps = append(beatmapsetBeatmaps, database.Beatmap{
				BeatmapId:     int32(beatmapId),
				BeatmapsetId:  currentBeatmapset.BeatmapsetId,
				CreatorId:     currentBeatmapset.CreatorId,
				Filename:      foundFilename,
				BeatmapMd5:    beatmapInfo.FileMd5,
				Version:       beatmapInfo.Version,
				TotalLength:   int32(totalLength),
				DrainTime:     int32(drainTime),
				CountObjects:  int32(countObjects),
				CountNormal:   int32(countNormal),
				CountSlider:   int32(countSliders),
				CountSpinner:  int32(countSpinners),
				DiffHp:        int8(diffHp),
				DiffCs:        int8(diffCs),
				DiffOd:        int8(diffOd),
				DiffStars:     -1,
				Playmode:      int8(playmode),
				RankingStatus: int8(rankingStatus),
				LastUpdate:    beatmapInfo.LastUpdate,
				ApproveDate:   beatmapInfo.ApprovedDate,
				SubmitDate:    beatmapInfo.SubmitDate,
				BeatmapSource: 0,
			})
		}

		mapsetVersions := ""

		for _, beatmapsetBeatmap := range beatmapsetBeatmaps {
			mapsetVersions += strconv.FormatInt(int64(beatmapsetBeatmap.BeatmapId), 10) + ","

			beatmapInsert, beatmapInsertErr := database.Database.Query("INSERT INTO waffle.beatmaps (beatmap_id, beatmapset_id, creator_id, filename, beatmap_md5, version, total_length, drain_time, count_objects, count_normal, count_slider, count_spinner, diff_hp, diff_cs, diff_od, diff_stars, playmode, ranking_status, last_update, submit_date, approve_date, beatmap_source) VALUEs (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", beatmapsetBeatmap.BeatmapId, beatmapsetBeatmap.BeatmapsetId, beatmapsetBeatmap.CreatorId, beatmapsetBeatmap.Filename, beatmapsetBeatmap.BeatmapMd5, beatmapsetBeatmap.Version, beatmapsetBeatmap.TotalLength, beatmapsetBeatmap.DrainTime, beatmapsetBeatmap.CountObjects, beatmapsetBeatmap.CountNormal, beatmapsetBeatmap.CountSlider, beatmapsetBeatmap.CountSpinner, beatmapsetBeatmap.DiffHp, beatmapsetBeatmap.DiffCs, beatmapsetBeatmap.DiffOd, beatmapsetBeatmap.DiffStars, beatmapsetBeatmap.Playmode, beatmapsetBeatmap.RankingStatus, beatmapsetBeatmap.LastUpdate, beatmapsetBeatmap.SubmitDate, beatmapsetBeatmap.ApproveDate, beatmapsetBeatmap.BeatmapSource)

			if beatmapInsert != nil {
				beatmapInsert.Close()
			}

			if beatmapInsertErr != nil {
				errors++
				logger.Printf("Failed to insert Beatmap ID %d into the database", beatmapsetBeatmap.BeatmapId)
			} else {
				logger.Printf("Inserted Beatmap of Beatmap ID %d into the database", beatmapsetBeatmap.BeatmapId)
			}
		}

		beatmapsetInsert, beatmapsetInsertErr := database.Database.Query("INSERT INTO waffle.beatmapsets (beatmapset_id, creator_id, artist, title, creator, source, tags, has_video, has_storyboard, bpm) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", currentBeatmapset.BeatmapsetId, currentBeatmapset.CreatorId, currentBeatmapset.Artist, currentBeatmapset.Title, currentBeatmapset.Creator, currentBeatmapset.Source, currentBeatmapset.Tags, currentBeatmapset.HasVideo, currentBeatmapset.HasStoryboard, currentBeatmapset.Bpm)

		if beatmapsetInsert != nil {
			beatmapsetInsert.Close()
		}

		if beatmapsetInsertErr != nil {
			errors++
			logger.Printf("Failed to insert Set ID %s into the database", setId)
		}

		timeTaken := time.Since(startTime)
		timeTaken = timeTaken

		if errors == 0 {
			logger.Printf("Successfully Processed Set ID %s! Took %dms", setId, timeTaken.Milliseconds())
		} else {
			logger.Printf("Processed Set ID %s with %d errors... Took %dms", setId, errors, timeTaken.Milliseconds())
		}

		//To make sure there's only 5 requests a second, yes peppy allows you to do 60 but stillill
		time.Sleep(200)
	}
}
