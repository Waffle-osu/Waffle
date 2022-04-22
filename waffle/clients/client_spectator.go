package clients

import (
	"Waffle/waffle/client_manager"
	"Waffle/waffle/packets"
)

func (client *Client) BroadcastToSpectators(packetFunction func(chan packets.BanchoPacket)) {
	client.spectatorMutex.Lock()

	for _, spectator := range client.spectators {
		packetFunction(spectator.GetPacketQueue())
	}

	client.spectatorMutex.Unlock()
}

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

func (client *Client) InformSpectatorLeft(spectatingClient client_manager.OsuClient) {
	client.spectatorMutex.Lock()

	delete(client.spectators, spectatingClient.GetUserId())

	packets.BanchoSendSpectatorLeft(client.PacketQueue, spectatingClient.GetUserId())

	client.spectatorMutex.Unlock()

	client.BroadcastToSpectators(func(packetQueue chan packets.BanchoPacket) {
		packets.BanchoSendFellowSpectatorLeft(packetQueue, spectatingClient.GetUserId())
	})
}

func (client *Client) InformSpectatorCantSpectate(spectateClient client_manager.OsuClient) {
	client.BroadcastToSpectators(func(packetQueue chan packets.BanchoPacket) {
		packets.BanchoSendSpectatorCantSpectate(packetQueue, spectateClient.GetUserId())
	})
}
