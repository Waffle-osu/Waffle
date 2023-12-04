package b1815

import (
	"Waffle/bancho/lobby"
	"time"
)

func (client *Client) GetIdleTimes() (lastRecieve time.Time, logon time.Time) {
	return client.lastReceive, client.logonTime
}

func (client *Client) GetFormattedJoinedChannels() string {
	channelString := ""

	for _, value := range client.joinedChannels {
		if value.ReadPrivileges == 0 {
			channelString += value.Name + " "
		}
	}

	return channelString
}

func (client *Client) GetMultiplayerLobby() *lobby.MultiplayerLobby {
	return client.currentMultiLobby
}
