package clients

import (
	"Waffle/waffle/lobby"
	"Waffle/waffle/packets"
)

func (client *Client) LeaveCurrentMatch() {
	if client.currentMultiLobby != nil {
		client.currentMultiLobby.Part(client)
		client.currentMultiLobby = nil
	}
}

func (client *Client) JoinMatch(match *lobby.MultiplayerLobby, password string) {
	client.LeaveCurrentMatch()

	if match.Join(client, password) {
		client.currentMultiLobby = match

		packets.BanchoSendMatchJoinSuccess(client.PacketQueue, match.MatchInformation)
		packets.BanchoSendChannelAvailableAutojoin(client.PacketQueue, "#multiplayer")
	} else {
		packets.BanchoSendMatchJoinFail(client.PacketQueue)
	}
}

func (client *Client) GetStatus() packets.StatusUpdate {
	return client.Status
}
