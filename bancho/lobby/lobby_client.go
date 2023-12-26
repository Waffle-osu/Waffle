package lobby

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/database"
)

// LobbyClient defines an Interface of what we need from a client to be able to do Multiplayer
type LobbyClient interface {
	// Retrieves the User ID of the current client
	GetUserId() int32
	// Gets the client's User Information
	GetUserData() database.User
	// Retrieves the Users privileges
	GetUserPrivileges() int32
	//Sends the equivilant of a chat message to this client
	SendChatMessage(sender string, content string, channel string)
	// This function is used to retrieve the client's Privileges
	GetUsername() string
	// Retrieves the Relevant User stats of this client, relevant meaning for the currently active mode.
	GetRelevantUserStats() database.UserStats
	// Gets the client's current Status
	GetUserStatus() base_packet_structures.StatusUpdate

	// Makes this client leave the match the player is currently in
	LeaveCurrentMatch()
	// Attempts to join the client to a multiplayer lobby
	JoinMatch(match *MultiplayerLobby, password string)
	// Gets the users away message, empty if none
	GetAwayMessage() string

	// Sends the equivilant of a Channel Join information/message to this client
	InformChannelJoin(chatClient chat.ChatClient, channel *chat.Channel)
	// Sends the equivilant of a Channel Part information/message to this client
	InformChannelPart(chatClient chat.ChatClient, channel *chat.Channel)

	//Sends the equivilant of a Lobby Join Packet
	BanchoLobbyJoin(userId int32)
	//Sends the equivilant of a Lobby Part Packet
	BanchoLobbyLeft(userId int32)

	// Informs the client of a new multiplayer match
	BanchoMatchNew(match base_packet_structures.MultiplayerMatch)
	// Informs the client of a update to a multiplayer match
	BanchoMatchUpdate(match base_packet_structures.MultiplayerMatch)
	// Informs the client that the match they're in has begun playing
	BanchoMatchStart(match base_packet_structures.MultiplayerMatch)
	// Informs the client that a multiplayer match has been disbanded
	BanchoMatchDisband(matchId int32)
	// Informs the client that they've been protomoted to host of the match
	BanchoMatchTransferHost()
	// Informs the client that all the other players in the match have loaded, so they can begin playing
	BanchoMatchAllPlayersLoaded()
	// Informs the client that all the other players finished playing
	BanchoMatchComplete()
	// Informs the client that all the other players also skipped the beginning break
	BanchoMatchSkip()
	// Lets the client know that player in slot `slot` skipped
	BanchoMatchPlayerSkipped(slot int32)
	// Lets the client know that player in slot `slot` failed
	BanchoMatchPlayerFailed(slot int32)
	// Lets the client know of another players current score in the lobby
	BanchoMatchScoreUpdate(scoreFrame base_packet_structures.ScoreFrame)

	// Lets the client know of another players new statistics
	BanchoOsuUpdate(stats database.UserStats, update base_packet_structures.StatusUpdate)

	// Lets the client know that a channel doesn't exist anymore
	BanchoChannelRevoked(channel string)

	// Silences the client until `untilUnix`
	SetSilencedUntilUnix(untilUnix int64)
	// Retrieves until what time the client is silenced until
	GetSilencedUntilUnix() int64

	// Gets the multiplayer lobby the client is currently in
	GetMultiplayerLobby() *MultiplayerLobby
	// Forcefully assigns a multiplayer lobby to a client
	AssignMultiplayerLobby(lobby *MultiplayerLobby)
	// Force joins a client into a chat channel
	AddJoinedChannel(channel *chat.Channel)
}
