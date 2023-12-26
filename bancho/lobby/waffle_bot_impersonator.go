package lobby

import (
	"Waffle/bancho/chat"
)

type LobbyWaffleBot struct{}

// This function is used to retrieve the client's Privileges
func (LobbyWaffleBot) GetUserPrivileges() int32 {
	return 31
}

// Sends the equivilant of a Chat message to the client
func (LobbyWaffleBot) SendChatMessage(sender string, content string, channel string) {

}

// Retrieves the Username of the current client
func (LobbyWaffleBot) GetUsername() string {
	return "WaffleBot"
}

// Retrieves the User ID of the current client
func (LobbyWaffleBot) GetUserId() int32 {
	return 1
}

// Retrieves the Away message of the client, empty if none.
func (LobbyWaffleBot) GetAwayMessage() string {
	return ""
}

// Sends the equivilant of a Channel Join information/message to this client
func (LobbyWaffleBot) InformChannelJoin(chatClient chat.ChatClient, channel *chat.Channel) {

}

// Sends the equivilant of a Channel Part information/message to this client
func (LobbyWaffleBot) InformChannelPart(chatClient chat.ChatClient, channel *chat.Channel) {

}

// Silences the client until `untilUnix`
func (LobbyWaffleBot) SetSilencedUntilUnix(untilUnix int64) {

}

// Retrieves until what time the client is silenced until
func (LobbyWaffleBot) GetSilencedUntilUnix() int64 {
	return 0
}
