package clients

import (
	"Waffle/waffle/packets"
)

func (client *Client) IsOfAdminPrivileges() bool {
	return client.UserData.Privileges&16 > 0
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
