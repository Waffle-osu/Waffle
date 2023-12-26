package client_manager

import (
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/database"
	"time"
)

// WaffleClient defines an Interface of what we need from client.Client to be able to manage this client in the ClientLists
type WaffleClient interface {
	// Retrieves this client's User ID
	GetUserId() int32
	// Retrieves the Username of the current client
	GetUsername() string
	// Retrieves the Relevant User stats of this client, relevant meaning for the currently active mode.
	GetRelevantUserStats() database.UserStats
	// Gets the client's current Status
	GetUserStatus() base_packet_structures.StatusUpdate
	// Gets the client's User Information
	GetUserData() database.User
	// Retrieves the client's Timezone
	GetClientTimezone() int32
	// Retrieves the Idle times of the client, when the client last received a packet, and when they logged on
	GetIdleTimes() (lastReceive time.Time, logonTime time.Time)
	// Retrieves a IRC formatted joined channel list.
	GetFormattedJoinedChannels() string

	// Closes the client's connection, and also gives a reason.
	CleanupClient(reason string)
	// Cuts the client's connection abruptly.
	Cut()
	// Gets the current client's away message. Empty if none.
	GetAwayMessage() string

	// Sends the equivilant of a osu! client quit message.
	BanchoHandleOsuQuit(userId int32, username string)
	// Sends the equivilant of a IRC client quit message.
	BanchoHandleIrcQuit(username string)

	// Sends the equivilant of a chat message to the client.
	BanchoIrcMessage(message base_packet_structures.Message)

	// Sends the equivilant of a statistics update to the client.
	// Used to inform of score submissions of other clients, and a difference in stats
	BanchoOsuUpdate(stats database.UserStats, update base_packet_structures.StatusUpdate)
	// Sends the equivilant of a Presence update to the client.
	// This is used to inform the client that a client exists. used to be done by just the Stats update.
	BanchoPresence(user database.User, stats database.UserStats, timezone int32)

	// Sends the equivilant of a annoucement/notification to this client.
	// in osu! it shows up as a notification
	BanchoAnnounce(message string)
	// Used to get the attention of the client.
	BanchoGetAttention()

	// Silences the client until `untilUnix`
	SetSilencedUntilUnix(untilUnix int64)
	// Retrieves until what time the client is silenced.
	GetSilencedUntilUnix() int64
}
