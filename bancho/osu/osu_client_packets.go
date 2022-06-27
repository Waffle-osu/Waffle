package osu

import (
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/database"
)

type OsuClientPacketsInterface interface {
	BanchoAnnounce(annoucement string)
	BanchoBeatmapInfoReply(infoReply base_packet_structures.BeatmapInfoReply)
	BanchoChannelAvailableAutojoin(channel string)
	BanchoChannelAvailable(channel string)
	BanchoChannelJoinSuccess(channel string)
	BanchoChannelRevoked(channel string)

	BanchoFellowSpectatorJoined(userId int32)
	BanchoFellowSpectatorLeft(userId int32)
	BanchoSpectatorJoined(userId int32)
	BanchoSpectatorLeft(userId int32)

	BanchoFriendsList(friendsList []database.FriendEntry)
	BanchoGetAttention()
	BanchoPing()
	BanchoProtocolNegotiation(protocolVersion int32)

	BanchoHandleOsuQuit(userId int32)

	BanchoLobbyJoin(userId int32)
	BanchoLobbyLeft(userId int32)

	BanchoLoginPermissions(permissions int32)
	BanchoLoginReply(userId int32)

	BanchoMatchAllPlayersLoaded()
	BanchoMatchComplete()
	BanchoMatchDisband(matchId int32)
	BanchoMatchJoinFail()
	BanchoMatchJoinSuccess(match base_packet_structures.MultiplayerMatch)
	BanchoMatchNew(match base_packet_structures.MultiplayerMatch)
	BanchoMatchStart(match base_packet_structures.MultiplayerMatch)
	BanchoMatchUpdate(match base_packet_structures.MultiplayerMatch)
	BanchoMatchPlayerFailed(slotId int32)
	BanchoMatchPlayerSkipped(slotId int32)
	BanchoMatchPlayerScoreUpdate(scoreFrame base_packet_structures.ScoreFrame)
	BanchoMatchSkip()
	BanchoMatchTransferHost()

	BanchoOsuUpdate(user database.UserStats, status base_packet_structures.StatusUpdate)
	BanchoPresence(user database.User, stats database.UserStats, timezone int32)

	BanchoIrcMessage(message base_packet_structures.Message)

	BanchoSpectateFrames(frames base_packet_structures.SpectatorFrameBundle)
	BanchoSpectatorCantSpectate(userId int32)
}
