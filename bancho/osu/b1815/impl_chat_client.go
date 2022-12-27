package b1815

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/common"
	"Waffle/database"
)

// GetUserPrivileges returns the users privileges
func (client *Client) GetUserPrivileges() int32 {
	return client.UserData.Privileges
}

// SendChatMessage directly sends a chat message to the user
//TODO@(Furball): maybe remove this in favor of just getting the packet queue and sending it that way?
func (client *Client) SendChatMessage(sender string, content string, channel string) {
	client.BanchoIrcMessage(base_packet_structures.Message{
		Sender:  sender,
		Target:  channel,
		Message: content,
	})
}

// GetUsername gets the clients username
func (client *Client) GetUsername() string {
	return client.UserData.Username
}

// GetAwayMessage gets the away message the user has set, if any
func (client *Client) GetAwayMessage() string {
	return client.awayMessage
}

func (client *Client) InformChannelJoin(chatClient chat.ChatClient, channel *chat.Channel) {

}

func (client *Client) InformChannelPart(chatClient chat.ChatClient, channel *chat.Channel) {

}

func (client *Client) GetClientType() common.ClientType {
	return common.ClientTypeOsu1816
}

func (client *Client) GetSilencedUntilUnix() int64 {
	return int64(client.UserData.SilencedUntil)
}

func (client *Client) SetSilencedUntilUnix(untilUnix int64) {
	client.UserData.SilencedUntil = uint64(untilUnix)

	database.Database.Query("UPDATE waffle.users SET silenced_until = ? WHERE user_id = ?", untilUnix, client.UserData.UserID)
}
