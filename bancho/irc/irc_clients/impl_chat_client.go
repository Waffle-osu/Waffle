package irc_clients

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/irc/irc_messages"
	"Waffle/database"
)

// This function is used to retrieve the client's Privileges
func (client *IrcClient) GetUserPrivileges() int32 {
	return client.UserData.Privileges
}

// Sends the equivilant of a Chat message to the client
func (client *IrcClient) SendChatMessage(sender string, content string, channel string) {
	client.packetQueue <- irc_messages.IrcSendPrivMsg(sender, channel, content)
}

// Retrieves the Username of the current client
func (client *IrcClient) GetUsername() string {
	return client.Username
}

// Retrieves the User ID of the current client
func (client *IrcClient) GetUserId() int32 {
	return int32(client.UserData.UserID)
}

// Retrieves the Away message of the client, empty if none.
func (client *IrcClient) GetAwayMessage() string {
	return client.awayMessage
}

// Sends the equivilant of a Channel Join information/message to this client
func (client *IrcClient) InformChannelJoin(chatClient chat.ChatClient, channel *chat.Channel) {
	client.packetQueue <- irc_messages.IrcSendJoin(chatClient.GetUsername(), channel.Name)
}

// Sends the equivilant of a Channel Part information/message to this client
func (client *IrcClient) InformChannelPart(chatClient chat.ChatClient, channel *chat.Channel) {
	client.packetQueue <- irc_messages.IrcSendPart(chatClient.GetUsername(), channel.Name)
}

// Retrieves until what time the client is silenced until
func (client *IrcClient) GetSilencedUntilUnix() int64 {
	return int64(client.UserData.SilencedUntil)
}

// Silences the client until `untilUnix`
func (client *IrcClient) SetSilencedUntilUnix(untilUnix int64) {
	client.UserData.SilencedUntil = uint64(untilUnix)

	database.Database.Query("UPDATE waffle.users SET silenced_until = ? WHERE user_id = ?", untilUnix, client.UserData.UserID)
}
