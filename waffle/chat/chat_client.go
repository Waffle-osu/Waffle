package chat

type ChatClient interface {
	IsOfAdminPrivileges() bool
	SendChatMessage(sender string, content string, channel string)
	GetUsername() string
	GetUserId() int32
}
