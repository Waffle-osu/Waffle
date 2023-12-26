package irc_clients

import (
	"Waffle/bancho/irc/irc_messages"
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/database"
	"Waffle/helpers/serialization"
	"time"
)

// Retrieves the Relevant User stats of this client, relevant meaning for the currently active mode.
func (client *IrcClient) GetRelevantUserStats() database.UserStats {
	//If the ID is below 1, it gets recognized as a IRC client
	//Inside osu! clients, because inside b1815 BanchoPresence
	//we use the stats.UserID instead of user.UserID, this is exactly why.
	minusOne := int32(-1)

	return database.UserStats{
		UserID:         uint64(minusOne),
		Mode:           0,
		Rank:           0,
		RankedScore:    0,
		TotalScore:     0,
		Level:          0,
		Accuracy:       0,
		Playcount:      0,
		CountSSH:       0,
		CountSS:        0,
		CountSH:        0,
		CountS:         0,
		CountA:         0,
		CountB:         0,
		CountC:         0,
		CountD:         0,
		Hit300:         0,
		Hit100:         0,
		Hit50:          0,
		HitMiss:        0,
		HitGeki:        0,
		HitKatu:        0,
		ReplaysWatched: 0,
	}
}

// Gets the client's current Status
func (client *IrcClient) GetUserStatus() base_packet_structures.StatusUpdate {
	return base_packet_structures.StatusUpdate{
		Status:          serialization.OsuStatusUnknown,
		StatusText:      "on IRC",
		BeatmapChecksum: "No!",
		CurrentMods:     0,
		Playmode:        0,
		BeatmapId:       -1,
	}
}

// Gets the client's User Information
func (client *IrcClient) GetUserData() database.User {
	return client.UserData
}

// Retrieves the client's Timezone
func (client *IrcClient) GetClientTimezone() int32 {
	return 0
}

// Sends the equivilant of a osu! client quit message.
func (client *IrcClient) BanchoHandleOsuQuit(userId int32, username string) {
	client.BanchoHandleIrcQuit(username)
}

// Sends the equivilant of a IRC client quit message.
func (client *IrcClient) BanchoHandleIrcQuit(username string) {
	client.packetQueue <- irc_messages.Message{
		Command:  "QUIT",
		Trailing: "Leaving",
	}
}

// Sends the equivilant of a Spectator Join message.
// Used to build a Spectator List
func (client *IrcClient) BanchoSpectatorJoined(userId int32) {
	//We don't do anything here cuz no spectator over IRC
}

// Sends the equivilant of a Spectator Leave message.
// Used to build a Spectator List
func (client *IrcClient) BanchoSpectatorLeft(userId int32) {
	//We don't do anything here cuz no spectator over IRC
}

// Sends the equivilant of a Fellow Spectator Join message.
// Used to build a Spectator List
func (client *IrcClient) BanchoFellowSpectatorJoined(userId int32) {
	//We don't do anything here cuz no spectator over IRC
}

// Sends the equivilant of a Fellow Spectator Leave message.
// Used to build a Spectator List
func (client *IrcClient) BanchoFellowSpectatorLeft(userId int32) {
	//We don't do anything here cuz no spectator over IRC
}

// Sends the equivilant of a Spectator can't spectate message.
// in osu! there's a seperate list for Spectators that don't have the map.
func (client *IrcClient) BanchoSpectatorCantSpectate(userId int32) {
	//We don't do anything here cuz no spectator over IRC
}

// Sends the equivilant of Spectator Replay Frames to the client.
// This contains the next replay data of the client that this client is spectating
func (client *IrcClient) BanchoSpectateFrames(frameBundle base_packet_structures.SpectatorFrameBundle) {
	//We don't do anything here cuz no spectator over IRC
}

// Sends the equivilant of a chat message to the client.
func (client *IrcClient) BanchoIrcMessage(message base_packet_structures.Message) {
	client.packetQueue <- irc_messages.IrcSendPrivMsg(message.Sender, message.Target, message.Message)
}

// Sends the equivilant of a statistics update to the client.
// Used to inform of score submissions of other clients, and a difference in stats
func (client *IrcClient) BanchoOsuUpdate(stats database.UserStats, update base_packet_structures.StatusUpdate) {

}

// Sends the equivilant of a Presence update to the client.
// This is used to inform the client that a client exists. used to be done by just the Stats update.
func (client *IrcClient) BanchoPresence(user database.User, stats database.UserStats, timezone int32) {

}

// Retrieves the Idle times of the client, when the client last received a packet, and when they logged on
func (client *IrcClient) GetIdleTimes() (lastRecieve time.Time, logon time.Time) {
	return client.lastReceive, client.logonTime
}

// Retrieves a IRC formatted joined channel list.
func (client *IrcClient) GetFormattedJoinedChannels() string {
	channelString := ""

	for _, value := range client.joinedChannels {
		if value.ReadPrivileges == 0 {
			channelString += value.Name + " "
		}
	}

	return channelString
}

// Sends the equivilant of a annoucement/notification to this client.
// in osu! it shows up as a notification
func (ircClient *IrcClient) BanchoAnnounce(message string) {
	ircClient.BanchoIrcMessage(base_packet_structures.Message{
		Sender:  "WaffleBot",
		Target:  ircClient.Username,
		Message: message,
	})
}

// Used to get the attention of the client.
func (ircClient *IrcClient) BanchoGetAttention() {

}
