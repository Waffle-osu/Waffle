package b1815

import (
	"Waffle/bancho/osu/b1815/packets"
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/database"
)

// GetUserId gets the user's user id
func (client *Client) GetUserId() int32 {
	return int32(client.UserData.UserID)
}

// GetPacketQueue gets the user's current packet queue to which they can queue packets to
func (client *Client) GetPacketQueue() chan packets.BanchoPacket {
	return client.PacketQueue
}

// GetRelevantUserStats returns the stats depending on what game mode the user currently is playing on
func (client *Client) GetRelevantUserStats() database.UserStats {
	var stats database.UserStats

	switch client.Status.Playmode {
	case packets.OsuGamemodeOsu:
		stats = client.OsuStats
	case packets.OsuGamemodeTaiko:
		stats = client.TaikoStats
	case packets.OsuGamemodeCatch:
		stats = client.CatchStats
	case packets.OsuGamemodeMania:
		stats = client.ManiaStats
	}

	return stats
}

// GetUserStatus gets the users current status
func (client *Client) GetUserStatus() base_packet_structures.StatusUpdate {
	return client.Status
}

// GetUserData gets the users data
func (client *Client) GetUserData() database.User {
	return client.UserData
}

// GetClientTimezone returns the clients timezone
func (client *Client) GetClientTimezone() int32 {
	return client.ClientData.Timezone
}
