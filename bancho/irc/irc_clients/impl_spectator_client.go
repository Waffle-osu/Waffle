package irc_clients

import (
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/bancho/spectator"
	"fmt"
	"strings"
)

// BroadcastToSpectators broadcasts a packet to all the people spectating `client`
func (client *IrcClient) BroadcastToSpectators(packetFunction func(client spectator.SpectatorClient)) {
	if !client.IsOsu {
		return
	}

	client.spectatorMutex.Lock()

	for _, spectator := range client.spectators {
		packetFunction(spectator)
	}

	client.spectatorMutex.Unlock()
}

// Sends the equivilant of a Spectator Join message.
// Used to build a Spectator List
func (client *IrcClient) BanchoSpectatorJoined(userId int32) {
	if !client.IsOsu {
		return
	}

	spectatingClient := spectator.ClientManager.GetClientById(userId)

	client.spectatorMutex.Lock()
	client.spectators[spectatingClient.GetUserId()] = spectatingClient
	client.spectatorMutex.Unlock()

	client.BroadcastToSpectators(func(client spectator.SpectatorClient) {
		client.BanchoFellowSpectatorJoined(spectatingClient.GetUserId())
	})
}

// Sends the equivilant of a Spectator Leave message.
// Used to build a Spectator List
func (client *IrcClient) BanchoSpectatorLeft(userId int32) {
	if !client.IsOsu {
		return
	}

	spectatingClient := spectator.ClientManager.GetClientById(userId)

	client.spectatorMutex.Lock()
	client.spectators[spectatingClient.GetUserId()] = spectatingClient
	client.spectatorMutex.Unlock()

	client.BroadcastToSpectators(func(client spectator.SpectatorClient) {
		client.BanchoFellowSpectatorJoined(spectatingClient.GetUserId())
	})
}

// Sends the equivilant of a Fellow Spectator Join message.
// Used to build a Spectator List
func (client *IrcClient) BanchoFellowSpectatorJoined(userId int32) {

}

// Sends the equivilant of a Fellow Spectator Leave message.
// Used to build a Spectator List
func (client *IrcClient) BanchoFellowSpectatorLeft(userId int32) {

}

// Sends the equivilant of a Spectator can't spectate message.
// in osu! there's a seperate list for Spectators that don't have the map.
func (client *IrcClient) BanchoSpectatorCantSpectate(userId int32) {
	if !client.IsOsu {
		return
	}
}

func (client *IrcClient) sendChk(frameBundle base_packet_structures.SpectatorFrameBundle) {
	spectatingClientUsername := client.spectatingClient.GetUsername()
	spectatingClientStatus := client.spectatingClient.GetUserStatus()
	channelName := fmt.Sprintf("#%s-osu-r", spectatingClientUsername)

	perfectAsNum := "False"

	if frameBundle.ScoreFrame.Perfect {
		perfectAsNum = "True"
	}

	syncFormattedScore := fmt.Sprintf("0|%s|%d|%d|%d|%d|%d|%d|%d|%d|%s|%d|0|0|0", spectatingClientUsername, frameBundle.ScoreFrame.TotalScore, frameBundle.ScoreFrame.MaxCombo, frameBundle.ScoreFrame.Count50, frameBundle.ScoreFrame.Count100, frameBundle.ScoreFrame.Count300, frameBundle.ScoreFrame.CountMiss, frameBundle.ScoreFrame.CountKatu, frameBundle.ScoreFrame.CountGeki, perfectAsNum, spectatingClientStatus.CurrentMods)
	chkMessage := fmt.Sprintf("CHK+%s+%d+%d+%s", spectatingClientStatus.BeatmapChecksum, frameBundle.ScoreFrame.Time, frameBundle.ScoreFrame.MaxCombo, syncFormattedScore)

	//Sends a sync point
	client.SendChatMessage(spectatingClientUsername, chkMessage, channelName)
}

// Sends the equivilant of Spectator Replay Frames to the client.
// This contains the next replay data of the client that this client is spectating
func (client *IrcClient) BanchoSpectateFrames(frameBundle base_packet_structures.SpectatorFrameBundle) {
	if !client.IsOsu {
		return
	}

	spectatingClientUsername := client.spectatingClient.GetUsername()
	channelName := fmt.Sprintf("#%s-osu-r", spectatingClientUsername)

	switch frameBundle.ReplayAction {
	//Skip
	case 2:
		client.SendChatMessage(spectatingClientUsername, "SKIP", channelName)
	//Completion
	case 3:
		fallthrough
	//Fail
	case 4:
		client.SendChatMessage(spectatingClientUsername, "END", channelName)
	}

	client.sendChk(frameBundle)

	replayData := ""

	for _, frame := range frameBundle.Frames {
		replayData += fmt.Sprintf("%d|%f|%f|%d,", frame.Time, frame.MouseX, frame.MouseY, frame.ButtonState)
	}

	replayData = strings.TrimSuffix(replayData, ",")

	client.SendChatMessage(spectatingClientUsername, replayData, channelName)
}
