package lobby

import (
	"Waffle/bancho/osu/b1815/packets"
	"Waffle/bancho/osu/base_packet_structures"
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
	GetUserStatus() base_packet_structures.StatusUpdate

	LeaveCurrentMatch()
	JoinMatch(match *MultiplayerLobby, password string)
	GetAwayMessage() string
}
