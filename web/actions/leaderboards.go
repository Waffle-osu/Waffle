package actions

import (
	"Waffle/database"
	"errors"
	"fmt"
)

type GetLeaderboardsRequest struct {
	BeatmapChecksum string
	Filename        string
	Playmode        byte
	RankType        byte
	Mods            int32
	RequesterUserId int32

	GetScores        bool
	GetOffset        bool
	GetRequesterBest bool
}

type LeaderboardScore struct {
	ScoreId  int64
	Username string
	UserId   int64
	Score    int64

	Hit300  int32
	Hit100  int32
	Hit50   int32
	HitGeki int32
	HitKatu int32

	MaxCombo int32
	Perfect  bool
	Mods     int32
	Rank     int32
	Date     string
}

type GetLeaderboardsResponse struct {
	Error error

	HasOsz2Version   bool
	SubmissionStatus int8
	BeatmapId        int32
	BeatmapSetId     int32
	TotalScores      int32

	OnlineOffset int32
	DisplayTitle string
	OnlineRating float64

	Scores []LeaderboardScore
}

func GetLeaderboards(request GetLeaderboardsRequest) GetLeaderboardsResponse {
	response := GetLeaderboardsResponse{
		Error:          nil,
		HasOsz2Version: false,
	}

	beatmap := database.Beatmap{}
	beatmapQueryResult := int8(0)

	beatmapset := database.Beatmapset{}
	beatmapsetQueryResult := int8(0)

	//Early 2007 clients don't send a filename.
	if request.Filename != "" {
		beatmapQueryResult, beatmap = database.BeatmapsGetByFilename(request.Filename)
	} else {
		beatmapQueryResult, beatmap = database.BeatmapsGetByMd5(request.BeatmapChecksum)
	}

	switch beatmapQueryResult {
	//Query failed
	case -2:
		response.Error = errors.New("Beatmap Query failed.")

		fallthrough
	//Not found
	case -1:
		response.SubmissionStatus = -1
		return response
	}

	//Map needs updating
	if request.BeatmapChecksum != beatmap.BeatmapMd5 {
		response.SubmissionStatus = 1

		return response
	}

	beatmapsetQueryResult, beatmapset = database.BeatmapsetsGetBeatmapsetById(beatmap.BeatmapsetId)

	if beatmapsetQueryResult != 0 {
		response.Error = errors.New("Beatmapset Query failed")

		return response
	}

	response.BeatmapId = beatmap.BeatmapId
	response.BeatmapSetId = beatmap.BeatmapsetId
	response.DisplayTitle = fmt.Sprintf("[bold:0,size:20]%s|%s\n", beatmapset.Artist, beatmapset.Title)

	switch beatmap.RankingStatus {
	//Unranked
	case 0:
		response.SubmissionStatus = 0

		return response
	//Ranked
	case 1:
		response.SubmissionStatus = 2
	//Approved
	case 2:
		response.SubmissionStatus = 3
	}

	if request.GetOffset {
		offsetQueryResult, offset := database.BeatmapOffsetsGetBeatmapOffset(beatmap.BeatmapId)

		if offsetQueryResult == 0 {
			response.OnlineOffset = offset.Offset
		}
	}

	return response
}
