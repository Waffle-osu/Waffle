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
