package client_manager

import (
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/database"
)

// WaffleClient defines an Interface of what we need from client.Client to be able to manage this client in the ClientLists
type WaffleClient interface {
	GetUserId() int32
	GetRelevantUserStats() database.UserStats
	GetUserStatus() base_packet_structures.StatusUpdate
	GetUserData() database.User
	GetClientTimezone() int32

	CleanupClient(reason string)
	Cut()
	GetAwayMessage() string

	BanchoHandleOsuQuit(userId int32)
}
