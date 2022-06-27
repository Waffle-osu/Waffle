package osu

import (
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/database"
)

type OsuClient interface {
	//client_manager.WaffleClient
	// {
	GetUserId() int32
	GetRelevantUserStats() database.UserStats
	GetUserStatus() base_packet_structures.StatusUpdate
	GetUserData() database.User
	GetClientTimezone() int32

	CleanupClient(reason string)
	Cut()
	GetAwayMessage() string
	// } client_manager.WaffleClient

	HandleBeatmapInfoRequest(infoRequest base_packet_structures.BeatmapInfoRequest)

	OsuClientPacketsInterface
}
