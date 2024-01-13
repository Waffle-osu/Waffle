package bss

import (
	"Waffle/database"

	"github.com/Waffle-osu/osu-parser/osu_parser"
)

type UploadTicket struct {
	Filename string
	Ticket   string
	Size     int64
	Metadata osu_parser.MetadataSection
	FileData []byte
}

type UploadRequest struct {
	UploadTickets map[string]UploadTicket
	HasVideo      bool
	HasStoryboard bool
	OszTicket     string
	BeatmapsetId  int64
	IsUpdate      bool

	Metadata osu_parser.MetadataSection
}

var uploadRequests map[int64]*UploadRequest = map[int64]*UploadRequest{}

func RegisterRequest(userId int64, uploadRequest *UploadRequest) (int64, error) {

	minBeatmapsetIdQuery := `
		SELECT final_beatmapset_id FROM (
			SELECT
				MIN(beatmapsets.beatmapset_id),
				CASE WHEN beatmapset_id IS NULL THEN 100000000 ELSE beatmapset_id END AS 'final_beatmapset_id'
			FROM 
				beatmapsets 
			WHERE beatmapsets.beatmapset_id > 100000000
		) a
	`

	query, queryErr := database.Database.Query(minBeatmapsetIdQuery)

	if queryErr != nil {
		return -1, queryErr
	}

	newSetId := int64(0)

	query.Next()
	scanErr := query.Scan(&newSetId)
	query.Close()

	if scanErr != nil {
		return -1, scanErr
	}

	uploadRequests[userId] = uploadRequest

	return newSetId, nil
}

func GetUploadRequest(userId int64) *UploadRequest {
	return uploadRequests[userId]
}

func DeleteUploadRequest(userId int64) {
	delete(uploadRequests, userId)
}
