package clients

import (
	"Waffle/bancho/lobby"
	"Waffle/bancho/packets"
)

// LeaveCurrentMatch makes the client leave the current match
func (client *Client) LeaveCurrentMatch() {
	if client.currentMultiLobby != nil {
		client.currentMultiLobby.Part(client)
		client.currentMultiLobby = nil
	}
}

// JoinMatch makes the client join a particular match
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
