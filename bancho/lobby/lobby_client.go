package lobby

import (
	"Waffle/bancho/packets"
	"Waffle/database"
)

type LobbyClient interface {
	GetPacketQueue() chan packets.BanchoPacket
	GetUserId() int32
	GetUserData() database.User
	GetUserPrivileges() int32
	SendChatMessage(sender string, content string, channel string)
	GetUsername() string
	GetRelevantUserStats() database.UserStats
	GetUserStatus() packets.StatusUpdate

	LeaveCurrentMatch()
	JoinMatch(match *MultiplayerLobby, password string)
	GetAwayMessage() string
}
