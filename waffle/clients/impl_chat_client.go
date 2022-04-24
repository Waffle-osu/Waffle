package clients

import (
	"Waffle/waffle/packets"
)

func (client *Client) GetUserPrivileges() int32 {
	return client.UserData.Privileges
}

func (client *Client) SendChatMessage(sender string, content string, channel string) {
	packets.BanchoSendIrcMessage(client.PacketQueue, packets.Message{Sender: sender, Message: content, Target: channel})
}

func (client *Client) GetUsername() string {
	return client.UserData.Username
}

func (client *Client) GetAwayMessage() string {
	return client.awayMessage
}
