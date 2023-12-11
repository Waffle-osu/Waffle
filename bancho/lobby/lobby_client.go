package lobby

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/common"
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

	InformChannelJoin(chatClient chat.ChatClient, channel *chat.Channel)
	InformChannelPart(chatClient chat.ChatClient, channel *chat.Channel)
	GetClientType() common.ClientType

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

	SetSilencedUntilUnix(untilUnix int64)
	GetSilencedUntilUnix() int64

	GetMultiplayerLobby() *MultiplayerLobby
	AssignMultiplayerLobby(lobby *MultiplayerLobby)
}
