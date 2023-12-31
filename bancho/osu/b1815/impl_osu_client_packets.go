package b1815

import (
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/bancho/spectator"
	"Waffle/database"
	"Waffle/helpers/packets"
	"bytes"
	"encoding/binary"
)

/*
	All these functions send out the packet
	their functions are named after, no
	real special functionality here.
*/

func (client *Client) BanchoAnnounce(annoucement string) {
	client.PacketQueue <- packets.Send(packets.BanchoAnnounce, annoucement)
}

func (client *Client) BanchoBeatmapInfoReply(infoReply base_packet_structures.BeatmapInfoReply) {
	client.PacketQueue <- packets.Send(packets.BanchoBeatmapInfoReply, infoReply)
}

func (client *Client) BanchoChannelAvailableAutojoin(channel string) {
	client.PacketQueue <- packets.Send(packets.BanchoChannelAvailableAutojoin, channel)
}

func (client *Client) BanchoChannelAvailable(channel string) {
	client.PacketQueue <- packets.Send(packets.BanchoChannelAvailable, channel)
}

func (client *Client) BanchoChannelJoinSuccess(channel string) {
	client.PacketQueue <- packets.Send(packets.BanchoChannelJoinSuccess, channel)
}

func (client *Client) BanchoChannelRevoked(channel string) {
	client.PacketQueue <- packets.Send(packets.BanchoChannelRevoked, channel)
}

func (client *Client) BanchoFellowSpectatorJoined(userId int32) {
	client.PacketQueue <- packets.Send(packets.BanchoFellowSpectatorJoined, userId)
}

func (client *Client) BanchoFellowSpectatorLeft(userId int32) {
	client.PacketQueue <- packets.Send(packets.BanchoFellowSpectatorLeft, userId)
}

func (client *Client) BanchoSpectatorJoined(userId int32) {
	client.PacketQueue <- packets.Send(packets.BanchoSpectatorJoined, userId)

	client.InformSpectatorJoin(spectator.ClientManager.GetClientById(userId))
}

func (client *Client) BanchoSpectatorLeft(userId int32) {
	client.PacketQueue <- packets.Send(packets.BanchoSpectatorLeft, userId)

	client.InformSpectatorLeft(spectator.ClientManager.GetClientById(userId))
}

func (client *Client) BanchoFriendsList(friendsList []database.FriendEntry) {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, int16(len(friendsList)))

	for _, friend := range friendsList {
		binary.Write(buf, binary.LittleEndian, int32(friend.User2))
	}

	client.PacketQueue <- packets.SendBytes(packets.BanchoFriendsList, buf.Bytes())
}

func (client *Client) BanchoGetAttention() {
	client.PacketQueue <- packets.SendEmpty(packets.BanchoGetAttention)
}

func (client *Client) BanchoPing() {
	client.PacketQueue <- packets.SendEmpty(packets.BanchoPing)
}

func (client *Client) BanchoProtocolNegotiation(protocolVersion int32) {
	client.PacketQueue <- packets.Send(packets.BanchoProtocolNegotiation, protocolVersion)
}

func (client *Client) BanchoHandleOsuQuit(userId int32, username string) {
	client.PacketQueue <- packets.Send(packets.BanchoHandleOsuQuit, userId)
}

func (client *Client) BanchoHandleIrcQuit(username string) {
	client.PacketQueue <- packets.Send(packets.BanchoHandleIrcQuit, username)
}

func (client *Client) BanchoLobbyJoin(userId int32) {
	client.PacketQueue <- packets.Send(packets.BanchoLobbyJoin, userId)
}

func (client *Client) BanchoLobbyLeft(userId int32) {
	client.PacketQueue <- packets.Send(packets.BanchoLobbyPart, userId)
}

func (client *Client) BanchoLoginPermissions(permissions int32) {
	client.PacketQueue <- packets.Send(packets.BanchoLoginPermissions, permissions)
}

func (client *Client) BanchoLoginReply(userId int32) {
	client.PacketQueue <- packets.Send(packets.BanchoLoginReply, userId)
}

func (client *Client) BanchoMatchAllPlayersLoaded() {
	client.PacketQueue <- packets.SendEmpty(packets.BanchoMatchAllPlayersLoaded)
}

func (client *Client) BanchoMatchComplete() {
	client.PacketQueue <- packets.SendEmpty(packets.BanchoMatchComplete)
}

func (client *Client) BanchoMatchDisband(matchId int32) {
	client.PacketQueue <- packets.Send(packets.BanchoMatchDisband, matchId)
}

func (client *Client) BanchoMatchJoinFail() {
	client.PacketQueue <- packets.SendEmpty(packets.BanchoMatchJoinFail)
}

func (client *Client) BanchoMatchJoinSuccess(match base_packet_structures.MultiplayerMatch) {
	client.PacketQueue <- packets.Send(packets.BanchoMatchJoinSuccess, match)
}

func (client *Client) BanchoMatchNew(match base_packet_structures.MultiplayerMatch) {
	client.PacketQueue <- packets.Send(packets.BanchoMatchNew, match)
}

func (client *Client) BanchoMatchStart(match base_packet_structures.MultiplayerMatch) {
	client.PacketQueue <- packets.Send(packets.BanchoMatchStart, match)
}

func (client *Client) BanchoMatchUpdate(match base_packet_structures.MultiplayerMatch) {
	client.PacketQueue <- packets.Send(packets.BanchoMatchUpdate, match)
}

func (client *Client) BanchoMatchPlayerFailed(slotId int32) {
	client.PacketQueue <- packets.Send(packets.BanchoMatchPlayerFailed, slotId)
}

func (client *Client) BanchoMatchPlayerSkipped(slotId int32) {
	client.PacketQueue <- packets.Send(packets.BanchoMatchPlayerSkipped, slotId)
}

func (client *Client) BanchoMatchScoreUpdate(scoreFrame base_packet_structures.ScoreFrame) {
	client.PacketQueue <- packets.Send(packets.BanchoMatchScoreUpdate, scoreFrame)
}

func (client *Client) BanchoMatchSkip() {
	client.PacketQueue <- packets.SendEmpty(packets.BanchoMatchSkip)
}

func (client *Client) BanchoMatchTransferHost() {
	client.PacketQueue <- packets.SendEmpty(packets.BanchoMatchTransferHost)
}

func (client *Client) BanchoOsuUpdate(user database.UserStats, status base_packet_structures.StatusUpdate) {
	stats := base_packet_structures.OsuStats{
		UserId:      int32(user.UserID),
		Status:      status,
		RankedScore: int64(user.RankedScore),
		Accuracy:    user.Accuracy,
		Playcount:   int32(user.Playcount),
		TotalScore:  int64(user.TotalScore),
		Rank:        int32(user.Rank),
	}

	client.PacketQueue <- packets.Send(packets.BanchoHandleOsuUpdate, stats)
}

func (client *Client) BanchoPresence(user database.User, stats database.UserStats, timezone int32) {
	//We're using stats.UserID instead of user.UserID for a small hack
	//Regarding Presence and IRC, because IRC client have a userid of -1 in the presence packet
	//this is how the osu client detects them, so IRC waffleclient implementation sends out a relevant
	//stats struct with the user id of -1 for precisely this purpose.
	presence := base_packet_structures.UserPresence{
		UserId:          int32(stats.UserID),
		Username:        user.Username,
		AvatarExtension: 1,
		Timezone:        uint8(timezone),
		Country:         uint8(user.Country),
		City:            "",
		Permissions:     uint8(user.Privileges),
		Longitude:       0,
		Latitude:        0,
		Rank:            int32(stats.Rank),
	}

	client.PacketQueue <- packets.Send(packets.BanchoUserPresence, presence)
}

func (client *Client) BanchoIrcMessage(message base_packet_structures.Message) {
	client.PacketQueue <- packets.Send(packets.BanchoSendMessage, message)
}

func (client *Client) BanchoSpectateFrames(frames base_packet_structures.SpectatorFrameBundle) {
	client.PacketQueue <- packets.Send(packets.BanchoSpectateFrames, frames)
}

func (client *Client) BanchoSpectatorCantSpectate(userId int32) {
	client.PacketQueue <- packets.Send(packets.BanchoSpectatorCantSpectate, userId)

	client.InformSpectatorCantSpectate(spectator.ClientManager.GetClientById(userId))
}
