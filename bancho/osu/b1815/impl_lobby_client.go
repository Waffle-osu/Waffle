package b1815

import (
	"Waffle/bancho/lobby"
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

		client.BanchoMatchJoinSuccess(match.MatchInformation)
		client.BanchoChannelAvailableAutojoin("#multiplayer")
	} else {
		client.BanchoMatchJoinFail()
	}
}

func (client *Client) AssignMultiplayerLobby(lobby *lobby.MultiplayerLobby) {
	client.currentMultiLobby = lobby
}
