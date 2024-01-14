package bss

import (
	"Waffle/database"
	"database/sql"

	"github.com/Waffle-osu/osu-parser/osu_parser"
)

func CheckBeatmapStatus(beatmapsetId int64, userData database.User, metadata *osu_parser.MetadataSection) (canEdit bool, exists bool, approved bool, setId int64, queryErrorOccured bool) {
	beatmapCountQuerySql := "SELECT COUNT(*), SUM(ranking_status) FROM beatmaps WHERE beatmapset_id = ?"

	beatmapCountQuery, beatmapCountQueryErr := database.Database.Query(beatmapCountQuerySql, beatmapsetId)

	if beatmapCountQueryErr != nil {
		return false, false, false, -1, true
	}

	count := int64(0)
	queryApproved := sql.NullInt64{}

	beatmapCountQuery.Next()
	scanErr := beatmapCountQuery.Scan(&count, &queryApproved)
	beatmapCountQuery.Close()

	if !queryApproved.Valid {
		queryApproved.Int64 = 0
	}

	if scanErr != nil {
		return false, false, false, -1, true
	}

	if count == 0 {
		return true, false, false, -1, false
	}

	//editable if not ranked/approved
	toReturnCanEdit := queryApproved.Int64 < count //TODO: && metadata.Creator == userData.Username

	return toReturnCanEdit, true, toReturnCanEdit, beatmapsetId, false
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
