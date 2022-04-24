package chat

type ChatClient interface {
	GetUserPrivileges() int32
	SendChatMessage(sender string, content string, channel string)
	GetUsername() string
	GetUserId() int32
	GetAwayMessage() string
}
