package b1815

import (
	"Waffle/bancho/osu"
	"Waffle/bancho/osu/b1815/packets"
)

// BroadcastToSpectators broadcasts a packet to all the people spectating `client`
func (client *Client) BroadcastToSpectators(packetFunction func(client osu.OsuClient)) {
	client.spectatorMutex.Lock()

	for _, spectator := range client.spectators {
		packetFunction(spectator)
	}

	client.spectatorMutex.Unlock()
}

// InformSpectatorJoin is called by a new spectator, informing this client that its now being watched
func (client *Client) InformSpectatorJoin(spectatingClient osu.OsuClient) {
	client.spectatorMutex.Lock()

	client.spectators[spectatingClient.GetUserId()] = spectatingClient

	for _, spectator := range client.spectators {
		spectator.BanchoFellowSpectatorJoined(spectator.GetUserId())
	}

	packets.BanchoSendSpectatorJoin(client.PacketQueue, spectatingClient.GetUserId())

	client.spectatorMutex.Unlock()

	client.BroadcastToSpectators(func(client osu.OsuClient) {
		client.BanchoFellowSpectatorJoined(spectatingClient.GetUserId())
	})
}

// InformSpectatorLeft is called by a spectator, informing that it has stopped watching
func (client *Client) InformSpectatorLeft(spectatingClient osu.OsuClient) {
	client.spectatorMutex.Lock()

	delete(client.spectators, spectatingClient.GetUserId())

	packets.BanchoSendSpectatorLeft(client.PacketQueue, spectatingClient.GetUserId())

	client.spectatorMutex.Unlock()

	client.BroadcastToSpectators(func(client osu.OsuClient) {
		client.BanchoFellowSpectatorLeft(spectatingClient.GetUserId())
	})
}

// InformSpectatorCantSpectate is called by a spectator, informing that it doesn't own the beatmap that is being played
func (client *Client) InformSpectatorCantSpectate(spectateClient osu.OsuClient) {
	client.BroadcastToSpectators(func(client osu.OsuClient) {
		client.BanchoSpectatorCantSpectate(spectateClient.GetUserId())
	})
}
