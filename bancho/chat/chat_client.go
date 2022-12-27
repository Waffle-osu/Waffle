package chat

import "Waffle/common"

// ChatClient defines an Interface of what we need from client.Client to be able to send messages
type ChatClient interface {
	GetUserPrivileges() int32
	SendChatMessage(sender string, content string, channel string)
	GetUsername() string
	GetUserId() int32
	GetAwayMessage() string
	InformChannelJoin(chatClient ChatClient, channel *Channel)
	InformChannelPart(chatClient ChatClient, channel *Channel)
	GetClientType() common.ClientType

	SetSilencedUntilUnix(untilUnix int64)
	GetSilencedUntilUnix() int64
}
