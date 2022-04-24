package lobby

import (
	"Waffle/waffle/database"
	"Waffle/waffle/packets"
)

type LobbyClient interface {
	GetPacketQueue() chan packets.BanchoPacket
	GetUserId() int32
	GetUserData() database.User
	IsOfAdminPrivileges() bool
	SendChatMessage(sender string, content string, channel string)
	GetUsername() string
	GetRelevantUserStats() database.UserStats
	GetStatus() packets.StatusUpdate

	LeaveCurrentMatch()
	JoinMatch(match *MultiplayerLobby, password string)
	GetAwayMessage() string
}
