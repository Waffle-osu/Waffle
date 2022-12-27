package irc_clients

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/irc/irc_messages"
	"Waffle/common"
)

func (client *IrcClient) GetUserPrivileges() int32 {
	return client.UserData.Privileges
}

func (client *IrcClient) SendChatMessage(sender string, content string, channel string) {
	client.packetQueue <- irc_messages.IrcSendPrivMsg(sender, channel, content)
}

func (client *IrcClient) GetUsername() string {
	return client.Username
}

func (client *IrcClient) GetUserId() int32 {
	return int32(client.UserData.UserID)
}

func (client *IrcClient) GetAwayMessage() string {
	return client.awayMessage
}

func (client *IrcClient) InformChannelJoin(chatClient chat.ChatClient, channel *chat.Channel) {
	client.packetQueue <- irc_messages.IrcSendJoin(chatClient.GetUsername(), channel.Name)
}

func (client *IrcClient) InformChannelPart(chatClient chat.ChatClient, channel *chat.Channel) {
	client.packetQueue <- irc_messages.IrcSendPart(chatClient.GetUsername(), channel.Name)
}

func (client *IrcClient) GetClientType() common.ClientType {
	return common.ClientTypeIrc
}

func (client *IrcClient) GetSilencedUntilUnix() int64 {
	return client.silencedUntil
}

func (client *IrcClient) SetSilencedUntilUnix(untilUnix int64) {
	client.silencedUntil = untilUnix
}
