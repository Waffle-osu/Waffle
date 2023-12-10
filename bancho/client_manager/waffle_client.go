package client_manager

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/lobby"
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/common"
	"Waffle/database"
	"time"
)

// WaffleClient defines an Interface of what we need from client.Client to be able to manage this client in the ClientLists
type WaffleClient interface {
	GetUserId() int32
	GetRelevantUserStats() database.UserStats
	GetUserStatus() base_packet_structures.StatusUpdate
	GetUserData() database.User
	GetClientTimezone() int32
	GetIdleTimes() (lastReceive time.Time, logonTime time.Time)
	GetFormattedJoinedChannels() string

	CleanupClient(reason string)
	Cut()
	GetAwayMessage() string

	BanchoHandleOsuQuit(userId int32)

	BanchoSpectatorJoined(userId int32)
	BanchoSpectatorLeft(userId int32)
	BanchoFellowSpectatorJoined(userId int32)
	BanchoFellowSpectatorLeft(userId int32)
	BanchoSpectatorCantSpectate(userId int32)
	BanchoSpectateFrames(frameBundle base_packet_structures.SpectatorFrameBundle)

	BanchoIrcMessage(message base_packet_structures.Message)

	BanchoOsuUpdate(stats database.UserStats, update base_packet_structures.StatusUpdate)
	BanchoPresence(user database.User, stats database.UserStats, timezone int32)

	BanchoAnnounce(message string)
	BanchoGetAttention()

	SetSilencedUntilUnix(untilUnix int64)
	GetSilencedUntilUnix() int64

	GetMultiplayerLobby() *lobby.MultiplayerLobby

	//heh fuckin entire LobbyClient here because of IRC Multiplayer Referree
	BanchoChannelRevoked(channel string)
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

	LeaveCurrentMatch()
	JoinMatch(match *lobby.MultiplayerLobby, password string)

	GetClientType() common.ClientType
	GetUserPrivileges() int32
	GetUsername() string

	InformChannelJoin(chatClient chat.ChatClient, channel *chat.Channel)
	InformChannelPart(chatClient chat.ChatClient, channel *chat.Channel)

	SendChatMessage(sender string, content string, channel string)
}
