package clients

func (client *Client) IsOfAdminPrivileges() bool {
	return client.UserData.Privileges&16 > 0
}

func (client *Client) SendChatMessage(sender string, content string, channel string) {

}

func (client *Client) GetUsername() string {
	return client.UserData.Username
}
