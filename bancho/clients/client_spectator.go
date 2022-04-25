package clients

import (
	"Waffle/bancho/client_manager"
	"Waffle/bancho/packets"
)

// BroadcastToSpectators broadcasts a packet to all the people spectating `client`
func (client *Client) BroadcastToSpectators(packetFunction func(chan packets.BanchoPacket)) {
	client.spectatorMutex.Lock()

	for _, spectator := range client.spectators {
		packetFunction(spectator.GetPacketQueue())
	}

	client.spectatorMutex.Unlock()
}

// InformSpectatorJoin is called by a new spectator, informing this client that its now being watched
func (client *Client) InformSpectatorJoin(spectatingClient client_manager.OsuClient) {
	client.spectatorMutex.Lock()

	client.spectators[spectatingClient.GetUserId()] = spectatingClient

	for _, spectator := range client.spectators {
		packets.BanchoSendFellowSpectatorJoined(spectatingClient.GetPacketQueue(), spectator.GetUserId())
	}

	packets.BanchoSendSpectatorJoin(client.PacketQueue, spectatingClient.GetUserId())

	client.spectatorMutex.Unlock()

	client.BroadcastToSpectators(func(packetQueue chan packets.BanchoPacket) {
		packets.BanchoSendFellowSpectatorJoined(packetQueue, spectatingClient.GetUserId())
	})
}

// InformSpectatorLeft is called by a spectator, informing that it has stopped watching
func (client *Client) InformSpectatorLeft(spectatingClient client_manager.OsuClient) {
	client.spectatorMutex.Lock()

	delete(client.spectators, spectatingClient.GetUserId())

	packets.BanchoSendSpectatorLeft(client.PacketQueue, spectatingClient.GetUserId())

	client.spectatorMutex.Unlock()

	client.BroadcastToSpectators(func(packetQueue chan packets.BanchoPacket) {
		packets.BanchoSendFellowSpectatorLeft(packetQueue, spectatingClient.GetUserId())
	})
}

// InformSpectatorCantSpectate is called by a spectator, informing that it doesn't own the beatmap that is being played
func (client *Client) InformSpectatorCantSpectate(spectateClient client_manager.OsuClient) {
	client.BroadcastToSpectators(func(packetQueue chan packets.BanchoPacket) {
		packets.BanchoSendSpectatorCantSpectate(packetQueue, spectateClient.GetUserId())
	})
}
