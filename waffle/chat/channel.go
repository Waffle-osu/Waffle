package chat

import (
	"sync"
)

type Channel struct {
	Name         string
	Description  string
	AdminChannel bool
	Clients      []ChatClient
	ClientMutex  sync.Mutex
}

func (channel *Channel) Join(client ChatClient) bool {
	if channel.AdminChannel && client.IsOfAdminPrivileges() == false {
		return false
	}

	channel.ClientMutex.Lock()
	channel.Clients = append(channel.Clients, client)
	channel.ClientMutex.Unlock()

	return true
}

func (channel *Channel) Leave(client ChatClient) {
	channel.ClientMutex.Lock()

	for index, value := range channel.Clients {
		if value == client {
			channel.Clients = append(channel.Clients[0:index], channel.Clients[index+1:]...)
		}
	}

	channel.ClientMutex.Unlock()
}

func (channel *Channel) SendMessage(sendingClient ChatClient, message string, target string) {
	channel.ClientMutex.Lock()

	for _, value := range channel.Clients {
		value.SendChatMessage(sendingClient.GetUsername(), message, target)
	}

	channel.ClientMutex.Unlock()
}
