package chat

import (
	"sync"
)

type Channel struct {
	Name            string
	Description     string
	ReadPrivileges  int32
	WritePrivileges int32
	Autojoin        bool
	Clients         []ChatClient
	ClientMutex     sync.Mutex
}

// Join Makes `client` Join the Channel, returns whether the attempt was successful
func (channel *Channel) Join(client ChatClient) bool {
	//If the user doesn't have read privileges, they shouldn't be allowed to join
	if (channel.ReadPrivileges & client.GetUserPrivileges()) <= 0 {
		return false
	}

	channel.ClientMutex.Lock()

	for _, chatUser := range channel.Clients {
		//Check for duplicate client
		if chatUser.GetUserId() == client.GetUserId() {
			channel.ClientMutex.Unlock()
			return true
		}
	}

	channel.Clients = append(channel.Clients, client)
	channel.ClientMutex.Unlock()

	// client.InformChannelJoin(client, channel)

	for _, chatUser := range channel.Clients {
		chatUser.InformChannelJoin(client, channel)
	}

	return true
}

// Leave Makes `client` Leave the Channel
func (channel *Channel) Leave(client ChatClient) {
	channel.ClientMutex.Lock()

	//Removes user from the Client list
	for index, value := range channel.Clients {
		if value == client {
			channel.Clients = append(channel.Clients[0:index], channel.Clients[index+1:]...)
		}
	}

	channel.ClientMutex.Unlock()

	for _, chatUser := range channel.Clients {
		chatUser.InformChannelPart(client, channel)
	}

	client.InformChannelPart(client, channel)
}

// SendMessage sends a message to the channel, `sendingClient` is the sender
func (channel *Channel) SendMessage(sendingClient ChatClient, message string, target string) {
	//If the user doesn't have write privileges, don't allow message sending
	if (channel.WritePrivileges & sendingClient.GetUserPrivileges()) <= 0 {
		sendingClient.SendChatMessage("WaffleBot", "You're not allowed to post in this channel! Your message has been discarded.", target)
		return
	}

	channel.ClientMutex.Lock()

	//Broadcast message to everyone in the channel
	for _, client := range channel.Clients {
		if client != sendingClient {
			client.SendChatMessage(sendingClient.GetUsername(), message, target)
		}
	}

	channel.ClientMutex.Unlock()
}
