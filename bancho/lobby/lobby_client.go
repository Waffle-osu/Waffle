package lobby

import (
	"Waffle/bancho/database"
	"Waffle/bancho/packets"
)

type LobbyClient interface {
	GetPacketQueue() chan packets.BanchoPacket
	GetUserId() int32
	GetUserData() database.User
	GetUserPrivileges() int32
	SendChatMessage(sender string, content string, channel string)
	GetUsername() string
	GetRelevantUserStats() database.UserStats
	GetStatus() packets.StatusUpdate

	LeaveCurrentMatch()
	JoinMatch(match *MultiplayerLobby, password string)
	GetAwayMessage() string
}
