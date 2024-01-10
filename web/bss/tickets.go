package bss

import "github.com/Waffle-osu/osu-parser/osu_parser"

type UploadTicket struct {
	Filename string
	Ticket   string
}

type UploadRequest struct {
	UploadTickets []UploadTicket
	HasVideo      bool
	HasStoryboard bool
	OszTicket     string

	Metadata osu_parser.MetadataSection
}

var uploadRequests map[int64]*UploadRequest = map[int64]*UploadRequest{}

func RegisterRequest(userId int64, uploadRequest UploadRequest) {
	uploadRequests[userId] = &uploadRequest
}

func GetUploadRequest(userId int64) *UploadRequest {
	return uploadRequests[userId]
}
