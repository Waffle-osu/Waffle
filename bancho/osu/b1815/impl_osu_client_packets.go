package b1815

import (
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/database"
	"Waffle/helpers/serialization"
	"bytes"
	"encoding/binary"
)

func (client *Client) BanchoAnnounce(annoucement string) {
	client.PacketQueue <- serialization.SendSerializableString(serialization.BanchoAnnounce, annoucement)
}

func (client *Client) BanchoBeatmapInfoReply(infoReply base_packet_structures.BeatmapInfoReply) {
	client.PacketQueue <- serialization.SendSerializable(serialization.BanchoBeatmapInfoReply, infoReply)
}

func (client *Client) BanchoChannelAvailableAutojoin(channel string) {
	client.PacketQueue <- serialization.SendSerializableString(serialization.BanchoChannelAvailableAutojoin, channel)
}

func (client *Client) BanchoChannelAvailable(channel string) {
	client.PacketQueue <- serialization.SendSerializableString(serialization.BanchoChannelAvailable, channel)
}

func (client *Client) BanchoChannelJoinSuccess(channel string) {
	client.PacketQueue <- serialization.SendSerializableString(serialization.BanchoChannelJoinSuccess, channel)
}

func (client *Client) BanchoChannelRevoked(channel string) {
	client.PacketQueue <- serialization.SendSerializableString(serialization.BanchoChannelRevoked, channel)
}

func (client *Client) BanchoFellowSpectatorJoined(userId int32) {
	client.PacketQueue <- serialization.SendSerializableInt(serialization.BanchoFellowSpectatorJoined, userId)
}

func (client *Client) BanchoFellowSpectatorLeft(userId int32) {
	client.PacketQueue <- serialization.SendSerializableInt(serialization.BanchoFellowSpectatorLeft, userId)
}

func (client *Client) BanchoSpectatorJoined(userId int32) {
	client.PacketQueue <- serialization.SendSerializableInt(serialization.BanchoSpectatorJoined, userId)
}

func (client *Client) BanchoSpectatorLeft(userId int32) {
	client.PacketQueue <- serialization.SendSerializableInt(serialization.BanchoSpectatorLeft, userId)
}

func (client *Client) BanchoFriendsList(friendsList []database.FriendEntry) {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, int16(len(friendsList)))

	for _, friend := range friendsList {
		binary.Write(buf, binary.LittleEndian, int32(friend.User2))
	}

	client.PacketQueue <- serialization.SendSerializableBytes(serialization.BanchoFriendsList, buf.Bytes())
}

func (client *Client) BanchoGetAttention() {
	client.PacketQueue <- serialization.SendEmptySerializable(serialization.BanchoGetAttention)
}

func (client *Client) BanchoPing() {
	client.PacketQueue <- serialization.SendEmptySerializable(serialization.BanchoPing)
}

func (client *Client) BanchoProtocolNegotiation(protocolVersion int32) {
	client.PacketQueue <- serialization.SendSerializableInt(serialization.BanchoProtocolNegotiation, protocolVersion)
}

func (client *Client) BanchoHandleOsuQuit(userId int32) {
	client.PacketQueue <- serialization.SendSerializableInt(serialization.BanchoHandleOsuQuit, userId)
}

func (client *Client) BanchoLobbyJoin(userId int32) {
	client.PacketQueue <- serialization.SendSerializableInt(serialization.BanchoLobbyJoin, userId)
}

func (client *Client) BanchoLobbyLeft(userId int32) {
	client.PacketQueue <- serialization.SendSerializableInt(serialization.BanchoLobbyPart, userId)
}

func (client *Client) BanchoLoginPermissions(permissions int32) {
	client.PacketQueue <- serialization.SendSerializableInt(serialization.BanchoLoginPermissions, permissions)
}

func (client *Client) BanchoLoginReply(userId int32) {
	client.PacketQueue <- serialization.SendSerializableInt(serialization.BanchoLoginReply, userId)
}

func (client *Client) BanchoMatchAllPlayersLoaded() {
	client.PacketQueue <- serialization.SendEmptySerializable(serialization.BanchoMatchAllPlayersLoaded)
}

func (client *Client) BanchoMatchComplete() {
	client.PacketQueue <- serialization.SendEmptySerializable(serialization.BanchoMatchComplete)
}

func (client *Client) BanchoMatchDisband(matchId int32) {
	client.PacketQueue <- serialization.SendSerializableInt(serialization.BanchoMatchDisband, matchId)
}

func (client *Client) BanchoMatchJoinFail() {
	client.PacketQueue <- serialization.SendEmptySerializable(serialization.BanchoMatchJoinFail)
}

func (client *Client) BanchoMatchJoinSuccess(match base_packet_structures.MultiplayerMatch) {
	client.PacketQueue <- serialization.SendSerializable(serialization.BanchoMatchJoinSuccess, match)
}

func (client *Client) BanchoMatchNew(match base_packet_structures.MultiplayerMatch) {
	client.PacketQueue <- serialization.SendSerializable(serialization.BanchoMatchNew, match)
}

func (client *Client) BanchoMatchStart(match base_packet_structures.MultiplayerMatch) {
	client.PacketQueue <- serialization.SendSerializable(serialization.BanchoMatchStart, match)
}

func (client *Client) BanchoMatchUpdate(match base_packet_structures.MultiplayerMatch) {
	client.PacketQueue <- serialization.SendSerializable(serialization.BanchoMatchUpdate, match)
}

func (client *Client) BanchoMatchPlayerFailed(slotId int32) {
	client.PacketQueue <- serialization.SendSerializableInt(serialization.BanchoMatchPlayerFailed, slotId)
}

func (client *Client) BanchoMatchPlayerSkipped(slotId int32) {
	client.PacketQueue <- serialization.SendSerializableInt(serialization.BanchoMatchPlayerSkipped, slotId)
}

func (client *Client) BanchoMatchPlayerScoreUpdate(scoreFrame base_packet_structures.ScoreFrame) {
	client.PacketQueue <- serialization.SendSerializable(serialization.BanchoMatchScoreUpdate, scoreFrame)
}

func (client *Client) BanchoMatchSkip() {
	client.PacketQueue <- serialization.SendEmptySerializable(serialization.BanchoMatchSkip)
}

func (client *Client) BanchoMatchTransferHost() {
	client.PacketQueue <- serialization.SendEmptySerializable(serialization.BanchoMatchTransferHost)
}

func (client *Client) BanchoOsuUpdate(user database.UserStats, status base_packet_structures.StatusUpdate) {
	client.PacketQueue <- serialization.SendSerializable(serialization.BanchoHandleOsuUpdate, status)
}

func (client *Client) BanchoPresence(user database.User, stats database.UserStats, timezone int32) {
	presence := base_packet_structures.UserPresence{
		UserId:          int32(user.UserID),
		Username:        user.Username,
		AvatarExtension: 0,
		Timezone:        uint8(timezone),
		Country:         uint8(user.Country),
		City:            "",
		Permissions:     uint8(user.Privileges),
		Longitude:       0,
		Latitude:        0,
		Rank:            int32(stats.Rank),
	}

	client.PacketQueue <- serialization.SendSerializable(serialization.BanchoUserPresence, presence)
}

func (client *Client) BanchoIrcMessage(message base_packet_structures.Message) {
	client.PacketQueue <- serialization.SendSerializable(serialization.BanchoSendMessage, message)
}

func (client *Client) BanchoSpectateFrames(frames base_packet_structures.SpectatorFrameBundle) {
	client.PacketQueue <- serialization.SendSerializable(serialization.BanchoSpectateFrames, frames)
}

func (client *Client) BanchoSpectatorCantSpectate(userId int32) {
	client.PacketQueue <- serialization.SendSerializableInt(serialization.BanchoSpectatorCantSpectate, userId)
}
