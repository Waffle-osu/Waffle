package b1815

import (
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/database"
	"Waffle/helpers/serialization"
)

// GetUserId gets the user's user id
func (client *Client) GetUserId() int32 {
	return int32(client.UserData.UserID)
}

// GetRelevantUserStats returns the stats depending on what game mode the user currently is playing on
func (client *Client) GetRelevantUserStats() database.UserStats {
	var stats database.UserStats

	switch client.Status.Playmode {
	case serialization.OsuGamemodeOsu:
		stats = client.OsuStats
	case serialization.OsuGamemodeTaiko:
		stats = client.TaikoStats
	case serialization.OsuGamemodeCatch:
		stats = client.CatchStats
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
