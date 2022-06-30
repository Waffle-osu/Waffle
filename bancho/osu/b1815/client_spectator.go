package b1815

import (
	"Waffle/bancho/client_manager"
)

// BroadcastToSpectators broadcasts a packet to all the people spectating `client`
func (client *Client) BroadcastToSpectators(packetFunction func(client client_manager.WaffleClient)) {
	client.spectatorMutex.Lock()

	for _, spectator := range client.spectators {
		packetFunction(spectator)
	}

	client.spectatorMutex.Unlock()
}

// InformSpectatorJoin is called by a new spectator, informing this client that its now being watched
func (client *Client) InformSpectatorJoin(spectatingClient client_manager.WaffleClient) {
	client.spectatorMutex.Lock()

	client.spectators[spectatingClient.GetUserId()] = spectatingClient

	for _, spectator := range client.spectators {
		spectator.BanchoFellowSpectatorJoined(spectator.GetUserId())
	}

	client.spectatorMutex.Unlock()

	client.BroadcastToSpectators(func(client client_manager.WaffleClient) {
		client.BanchoFellowSpectatorJoined(spectatingClient.GetUserId())
	})
}

// InformSpectatorLeft is called by a spectator, informing that it has stopped watching
func (client *Client) InformSpectatorLeft(spectatingClient client_manager.WaffleClient) {
	client.spectatorMutex.Lock()

	delete(client.spectators, spectatingClient.GetUserId())

	client.spectatorMutex.Unlock()

	client.BroadcastToSpectators(func(client client_manager.WaffleClient) {
		client.BanchoFellowSpectatorLeft(spectatingClient.GetUserId())
	})
}

// InformSpectatorCantSpectate is called by a spectator, informing that it doesn't own the beatmap that is being played
func (client *Client) InformSpectatorCantSpectate(spectatingClient client_manager.WaffleClient) {
	client.BroadcastToSpectators(func(client client_manager.WaffleClient) {
		client.BanchoSpectatorCantSpectate(spectatingClient.GetUserId())
	})
}
