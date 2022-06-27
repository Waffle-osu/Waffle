package lobby

import (
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/database"
	"Waffle/helpers/serialization"
)

type LobbyClient interface {
	GetUserId() int32
	GetUserData() database.User
	GetUserPrivileges() int32
	SendChatMessage(sender string, content string, channel string)
	GetUsername() string
	GetRelevantUserStats() database.UserStats
	GetUserStatus() base_packet_structures.StatusUpdate

	LeaveCurrentMatch()
	JoinMatch(match *MultiplayerLobby, password string)
	GetAwayMessage() string
}
