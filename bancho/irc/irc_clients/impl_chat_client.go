package irc_clients

func (client *IrcClient) GetUserPrivileges() int32 {
	return client.UserData.Privileges
}

func (client *IrcClient) SendChatMessage(sender string, content string, channel string) {
	
}
