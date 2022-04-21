package clients

func (client *Client) IsOfAdminPrivileges() bool {
	return client.UserData.Privileges&16 > 0
}
