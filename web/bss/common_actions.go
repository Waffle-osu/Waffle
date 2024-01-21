package bss

import (
	"Waffle/database"
	"Waffle/utils"
	"database/sql"
	"math"

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

	return insertBeatmapErr
}
