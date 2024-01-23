package bss

import (
	"Waffle/database"

	"github.com/Waffle-osu/osu-parser/osu_parser"
)

type UploadTicket struct {
	Filename  string
	Ticket    string
	Size      int64
	ParsedOsu osu_parser.OsuFile
	FileData  []byte
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

var uploadRequests map[int32]*UploadRequest = map[int32]*UploadRequest{}

func RegisterRequest(userId int32, uploadRequest *UploadRequest) (int64, error) {
	minBeatmapsetIdQuery := `
		SELECT final_beatmapset_id + 1 FROM (
			SELECT 
				next_id,
				CASE WHEN next_id IS NULL THEN (100000000-1) ELSE next_id END AS 'final_beatmapset_id'
			FROM (
				SELECT MAX(beatmapset_id) AS 'next_id' FROM beatmapsets WHERE beatmapset_id >= 100000000
			) a
		) b
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

func GetUploadRequest(userId int32) *UploadRequest {
	return uploadRequests[userId]
}

func DeleteUploadRequest(userId int32) {
	delete(uploadRequests, userId)
}
