package b1815

import (
	"Waffle/bancho/spectator"
)

// BroadcastToSpectators broadcasts a packet to all the people spectating `client`
func (client *Client) BroadcastToSpectators(packetFunction func(client spectator.SpectatorClient)) {
	client.spectatorMutex.Lock()

	for _, spectator := range client.spectators {
		packetFunction(spectator)
	}

	client.spectatorMutex.Unlock()
}

// InformSpectatorJoin is called by a new spectator, informing this client that its now being watched
func (client *Client) InformSpectatorJoin(spectatingClient spectator.SpectatorClient) {
	client.spectatorMutex.Lock()
	client.spectators[spectatingClient.GetUserId()] = spectatingClient
	client.spectatorMutex.Unlock()

	client.BroadcastToSpectators(func(client spectator.SpectatorClient) {
		client.BanchoFellowSpectatorJoined(spectatingClient.GetUserId())
	})
}

// InformSpectatorLeft is called by a spectator, informing that it has stopped watching
func (client *Client) InformSpectatorLeft(spectatingClient spectator.SpectatorClient) {
	client.spectatorMutex.Lock()
	delete(client.spectators, spectatingClient.GetUserId())
	client.spectatorMutex.Unlock()

	client.BroadcastToSpectators(func(client spectator.SpectatorClient) {
		client.BanchoFellowSpectatorLeft(spectatingClient.GetUserId())
	})
}

// InformSpectatorCantSpectate is called by a spectator, informing that it doesn't own the beatmap that is being played
func (client *Client) InformSpectatorCantSpectate(spectatingClient spectator.SpectatorClient) {
	client.BroadcastToSpectators(func(client spectator.SpectatorClient) {
		client.BanchoSpectatorCantSpectate(spectatingClient.GetUserId())
	})
}
