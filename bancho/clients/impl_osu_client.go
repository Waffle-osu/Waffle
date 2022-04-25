package clients

import (
	"Waffle/bancho/database"
	"Waffle/bancho/packets"
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
		break
	case packets.OsuGamemodeTaiko:
		stats = client.TaikoStats
		break
	case packets.OsuGamemodeCatch:
		stats = client.CatchStats
		break
	case packets.OsuGamemodeMania:
		stats = client.ManiaStats
		break
	}

	return stats
}

// GetUserStatus gets the users current status
//TODO@(Furball): Theres 2 functions that do the exact same thing...
func (client *Client) GetUserStatus() packets.StatusUpdate {
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
