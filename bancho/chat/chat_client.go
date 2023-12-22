package chat

import "Waffle/common"

// ChatClient defines an Interface of what we need from client.Client to be able to send messages
type ChatClient interface {
	// This function is used to retrieve the client's Privileges
	GetUserPrivileges() int32
	// Sends the equivilant of a Chat message to the client
	SendChatMessage(sender string, content string, channel string)
	// Retrieves the Username of the current client
	GetUsername() string
	// Retrieves the User ID of the current client
	GetUserId() int32
	// Retrieves the Away message of the client, empty if none.
	GetAwayMessage() string
	// Sends the equivilant of a Channel Join information/message to this client
	InformChannelJoin(chatClient ChatClient, channel *Channel)
	// Sends the equivilant of a Channel Part information/message to this client
	InformChannelPart(chatClient ChatClient, channel *Channel)
	// Gets what kind of client the client is.
	GetClientType() common.ClientType

	// Silences the client until `untilUnix`
	SetSilencedUntilUnix(untilUnix int64)
	// Retrieves until what time the client is silenced until
	GetSilencedUntilUnix() int64
}
