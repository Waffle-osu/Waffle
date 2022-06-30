package lobby

import (
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/database"
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

	BanchoLobbyJoin(userId int32)
	BanchoLobbyLeft(userId int32)

	BanchoMatchNew(match base_packet_structures.MultiplayerMatch)
	BanchoMatchUpdate(match base_packet_structures.MultiplayerMatch)
	BanchoMatchStart(match base_packet_structures.MultiplayerMatch)
	BanchoMatchDisband(matchId int32)
	BanchoMatchTransferHost()
	BanchoMatchAllPlayersLoaded()
	BanchoMatchComplete()
	BanchoMatchSkip()
	BanchoMatchPlayerSkipped(slot int32)
	BanchoMatchPlayerFailed(slot int32)
	BanchoMatchScoreUpdate(scoreFrame base_packet_structures.ScoreFrame)

	BanchoOsuUpdate(stats database.UserStats, update base_packet_structures.StatusUpdate)

	BanchoChannelRevoked(channel string)
}
