package clients

import (
	"Waffle/waffle/database"
	"Waffle/waffle/packets"
)

func (client *Client) GetUserId() int32 {
	return int32(client.UserData.UserID)
}

func (client *Client) GetPacketQueue() chan packets.BanchoPacket {
	return client.PacketQueue
}

func (client *Client) GetRelevantUserStats() database.UserStats {
	var stats database.UserStats

	switch client.Status.CurrentPlaymode {
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

func (client *Client) GetUserStatus() packets.OsuStatus {
	return client.Status
}

func (client *Client) GetUserData() database.User {
	return client.UserData
}

func (client *Client) GetClientTimezone() int32 {
	return client.ClientData.Timezone
}
