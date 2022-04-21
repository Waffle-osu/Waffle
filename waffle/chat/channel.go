package chat

import (
	"sync"
)

type Channel struct {
	Name         string
	Description  string
	AdminChannel bool
	Clients      []*ChatClient
	ClientMutex  sync.Mutex
}

func (channel Channel) Join(client ChatClient) bool {
	if channel.AdminChannel && client.IsOfAdminPrivileges() == false {
		return false
	}

	channel.ClientMutex.Lock()
	channel.Clients = append(channel.Clients, &client)
	channel.ClientMutex.Unlock()

	return true
}
