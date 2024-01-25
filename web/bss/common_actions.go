package bss

import (
	"Waffle/config"
	"Waffle/database"
	"Waffle/helpers"
	"Waffle/utils"
	"Waffle/web/bss/thumbnail"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Waffle-osu/osu-parser/osu_parser"
)

func getMapCountAndApprovedStatus(beatmapsetId int64) (count int64, approved bool, err error) {
	beatmapCountQuerySql := "SELECT COUNT(*), SUM(ranking_status) FROM beatmaps WHERE beatmapset_id = ?"

	beatmapCountQuery, beatmapCountQueryErr := database.Database.Query(beatmapCountQuerySql, beatmapsetId)

	if beatmapCountQueryErr != nil {
		return 0, false, beatmapCountQueryErr
	}

	queryCount := int64(0)
	queryApproved := sql.NullInt64{}

	beatmapCountQuery.Next()
	scanErr := beatmapCountQuery.Scan(&queryCount, &queryApproved)
	beatmapCountQuery.Close()

	if scanErr != nil {
		return 0, false, scanErr
	}

	if !queryApproved.Valid {
		queryApproved.Int64 = 0
	}

	approvedCalc := int64(math.Floor(float64(queryApproved.Int64) / float64(queryCount)))

	return queryCount, approvedCalc > count, nil
}

func CheckBeatmapStatus(beatmapsetId int64, userData database.User, metadata *osu_parser.MetadataSection) (canEdit bool, exists bool, approved bool, setId int64, queryErrorOccured bool) {
	count, queryApproved, err := getMapCountAndApprovedStatus(beatmapsetId)

	if err != nil {
		return false, false, false, -1, true
	}

	if count == 0 && metadata != nil {
		countSets := int64(0)
		negativeUserId := -int64(userData.UserID)
		//Try over metadata
		overMetadataSql := "SELECT COUNT(beatmapset_id) FROM beatmapsets WHERE artist = ? AND title = ? AND creator_id = ?"
		overMetadataQuery, overMetadataQueryErr := database.Database.Query(overMetadataSql, metadata.Artist, metadata.Title, negativeUserId)

		if overMetadataQueryErr != nil {
			return false, false, false, -1, true
		}

		overMetadataQuery.Next()
		setCountScanErr := overMetadataQuery.Scan(&countSets)
		overMetadataQuery.Close()

		if setCountScanErr != nil {
			return false, false, false, -1, true
		}

		if countSets == 0 {
			return true, false, false, -1, false
		}

		if countSets > 0 {
			//There shouldn't be more than one but sefjksdkfsbndlfbdsf
			getSetIdSql := "SELECT beatmapset_id FROM beatmapsets WHERE artist = ? AND title = ? AND creator_id = ? LIMIT 1"
			getSetIdQuery, getSetIdQueryErr := database.Database.Query(getSetIdSql, metadata.Artist, metadata.Title, negativeUserId)

			if getSetIdQueryErr != nil {
				return false, false, false, -1, true
			}

			foundSetId := int64(0)

			getSetIdQuery.Next()
			foundSetScanErr := getSetIdQuery.Scan(&foundSetId)
			getSetIdQuery.Close()

			if foundSetScanErr != nil {
				return false, false, false, -1, true
			}

			_, queryApproved, err := getMapCountAndApprovedStatus(beatmapsetId)

			if err != nil {
				return false, false, false, -1, true
			}

			return true, true, queryApproved, foundSetId, false
		}

		return true, false, false, -1, false
	}

	//editable if not ranked/approved
	toReturnCanEdit := !queryApproved && metadata.Creator == userData.Username

	return toReturnCanEdit, true, queryApproved, beatmapsetId, false
}

func GetNextBssBeatmapId() int64 {
	beatmapIdSql := `
		SELECT final_beatmap_id + 1 FROM (
			SELECT 
				next_id,
				CASE WHEN next_id IS NULL THEN (100000000-1) ELSE next_id END AS 'final_beatmap_id'
			FROM (
				SELECT MAX(beatmap_id) AS 'next_id' FROM beatmaps WHERE beatmap_id >= 100000000
			) a
		) b
	`

	beatmapIdQuery, beatmapIdErr := database.Database.Query(beatmapIdSql)
	result := int64(0)

	if beatmapIdErr != nil {
		return -1
	}

	beatmapIdQuery.Next()
	scanErr := beatmapIdQuery.Scan(&result)
	beatmapIdQuery.Close()

	if scanErr != nil {
		return -1
	}

	return result
}

func InsertIntoBeatmaps(file osu_parser.OsuFile, setId int64, userId int32, filename string) error {
	newBeatmapId := GetNextBssBeatmapId()

	minVersion := utils.VersionOsuFile(file)

	insertBeatmapSql := "INSERT INTO beatmaps (beatmap_id, beatmapset_id, creator_id, filename, beatmap_md5, version, total_length, drain_time, count_objects, count_normal, count_slider, count_spinner, diff_hp, diff_cs, diff_od, diff_stars, playmode, ranking_status, last_update, submit_date, approve_date, beatmap_source, status_valid_from_version, status_valid_to_version) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP(), CURRENT_TIMESTAMP(), '1000-01-01 00:00:00.000000', ?, ?, ?)"

	_, insertBeatmapErr :=
		database.Database.Query(
			insertBeatmapSql,
			newBeatmapId,
			setId,
			-userId,
			filename,
			file.Md5Hash,
			file.Metadata.Version,
			file.Length,
			file.DrainLength,
			len(file.HitObjects.List),
			file.HitObjects.CountNormal,
			file.HitObjects.CountSlider,
			file.HitObjects.CountSpinner,
			file.Difficulty.HPDrainRate,
			file.Difficulty.CircleSize,
			file.Difficulty.OverallDifficulty,
			-1,
			byte(file.General.Mode),
			0,
			1,
			minVersion,
			99999999,
		)

	if insertBeatmapErr == nil {
		go RunAndCreateDiffCalc(file, newBeatmapId, setId)
	}

	return insertBeatmapErr
}

func CreateTicket(user database.User, filename string, oszTicket bool) string {
	prefix := "osu"

	if oszTicket {
		prefix = "osz"
	}

	ticketFormat := fmt.Sprintf("%s::%d-%s-%s-%s", prefix, time.Now().Unix(), user.Username, filename, user.Password)
	ticketBytes := sha256.Sum256([]byte(ticketFormat))
	ticketHashed := ticketBytes[:]
	ticket := hex.EncodeToString(ticketHashed)

	return ticket
}

func GenerateThumbnail(uploadTicket UploadTicket, setId int64) {
	tempOszDir := fmt.Sprintf("bss_temp/oszs/%d", setId)

	//take its background and generate the thumbnail
	backgroundFilename := ""

	for _, event := range uploadTicket.ParsedOsu.Events.Events {
		if event.EventType == osu_parser.EventTypeBackground {
			backgroundFilename = event.BackgroundImage

			break
		}
	}

	backgroundFilename = strings.TrimPrefix(backgroundFilename, "\"")
	backgroundFilename = strings.TrimSuffix(backgroundFilename, "\"")

	imagePath := fmt.Sprintf("%s/%s", tempOszDir, backgroundFilename)

	generator := thumbnail.NewGenerator(thumbnail.Generator{
		Scaler: "CatmullRom",
	})

	generator.Width = 80
	generator.Height = 64

	image, imageErr := generator.NewImageFromFile(imagePath)

	if imageErr == nil {
		thumbnailBytes, thumbnailErr := generator.CreateThumbnail(image)

		if thumbnailErr == nil {
			os.WriteFile(fmt.Sprintf("direct_thumbnails/%d", setId), thumbnailBytes, 0644)
		} else {
			helpers.Logger.Printf("BSS:U %s failed to generate small thumbnail", uploadTicket.Filename)
		}

		generator.Width = 160
		generator.Height = 120

		thumbnailBytes, thumbnailErr = generator.CreateThumbnail(image)

		if thumbnailErr == nil {
			os.WriteFile(fmt.Sprintf("direct_thumbnails/%dl", setId), thumbnailBytes, 0644)
		} else {
			helpers.Logger.Printf("BSS:U %s failed to generate large thumbnail", uploadTicket.Filename)
		}
	} else {
		helpers.Logger.Printf("BSS:U %s failed to load image", uploadTicket.Filename)
	}
}

func CreateMp3Preview(audioFilename string, previewTimeMs int32, setId int64) {
	if config.FFMPEGPath == "" {
		helpers.Logger.Printf("FFMPEG Path not set, mp3 preview not generated.")

		return
	}

	tempOszDir := fmt.Sprintf("bss_temp/oszs/%d", setId)

	//And the mp3 preview
	previewStart := int(float64(previewTimeMs) / 1000.0)
	toString := strconv.FormatInt(int64(previewStart), 10)

	mp3Path := fmt.Sprintf("%s/%s", tempOszDir, audioFilename)
	outPath := fmt.Sprintf("mp3_previews/%d", setId)

	result, err := exec.Command(config.FFMPEGPath, "-ss", toString, "-t", "10", "-i", mp3Path, "-codec:a", "libmp3lame", "-b:a", "64k", outPath+".mp3").Output()

	if err != nil {
		if errors.Is(err, exec.ErrDot) {
			executable, execErr := os.Executable()

			if execErr == nil {
				ffmpegPath := fmt.Sprintf("%s/%s", filepath.Dir(executable), config.FFMPEGPath)

				result, err = exec.Command(ffmpegPath, "-ss", toString, "-t", "10", "-i", mp3Path, "-codec:a", "libmp3lame", "-b:a", "64k", outPath+".mp3").Output()
			}
		}

		if result != nil && err != nil {
			helpers.Logger.Printf("FFMPEG failed to generate mp3 preview: %s :::: %s", err.Error(), string(result))
		}
	}

	os.Rename(outPath+".mp3", outPath)
}
