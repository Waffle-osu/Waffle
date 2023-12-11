package lobby

import (
	"Waffle/bancho/osu/base_packet_structures"
	"fmt"
	"strings"
)

// WaffleBotCommandHelp !mp
func LobbyHandleCommandMultiplayer(sender LobbyClient, message string) []string {
	splitMessage := strings.Split(message, " ")

	if len(splitMessage) == 0 {
		return []string{
			"what",
		}
	}

	args := splitMessage[1:]

	subcommand := strings.ToLower(args[0])
	senderLobby := sender.GetMultiplayerLobby()

	if senderLobby == nil && subcommand != "make" {
		return []string{
			fmt.Sprintf("%s: Command only valid inside multiplayer lobbies!", subcommand),
		}
	}

	if len(args) == 0 {
		return []string{
			"Subcommand missing!",
		}
	}

	switch strings.ToLower(args[0]) {
	case "make":
		return MpCommandMake(sender, args)
	case "invite":
		return MpCommandInvite(sender, args)
	case "lock":
		return MpCommandLock(sender, args)
	case "unlock":
		return MpCommandUnlock(sender, args)
	case "size":
		return MpCommandSize(sender, args)
	case "set":
		return MpCommandSet(sender, args)
	case "move":
		return MpCommandMove(sender, args)
	case "team":
		return MpCommandTeam(sender, args)
	case "host":
		return MpCommandHost(sender, args)
	case "settings":
		return MpCommandSettings(sender, args)
	case "start":
		return MpCommandStart(sender, args)
	case "abort":
		return MpCommandAbort(sender, args)
	case "map":
		return MpCommandMap(sender, args)
	case "mods":
		return MpCommandMods(sender, args)
	case "timer":
		return MpCommandTimer(sender, args)
	case "aborttimer":
		return MpCommandAbortTimer(sender, args)
	case "kick":
		return MpCommandKick(sender, args)
	case "password":
		return MpCommandPassword(sender, args)
	case "close":
		return MpCommandClose(sender, args)
	}

	return []string{
		"!mp: Unknown Subcommand!",
	}
}

func MpCommandMake(sender LobbyClient, args []string) []string {
	if len(args) < 2 {
		return []string{
			"!mp make: Lobby name required!",
		}
	}

	lobbyName := ""

	for i := 1; i != len(args); i++ {
		lobbyName += args[i]
	}

	newLobby := CreateNewMultiMatch(base_packet_structures.MultiplayerMatch{
		MatchId:          0,
		InProgress:       false,
		MatchType:        base_packet_structures.MultiplayerMatchTypeHeadToHead,
		ActiveMods:       0,
		GameName:         lobbyName,
		GamePassword:     "",
		BeatmapName:      "No map selected.",
		BeatmapId:        1,
		BeatmapChecksum:  "",
		HostId:           sender.GetUserId(),
		Playmode:         0,
		MatchScoringType: base_packet_structures.MultiplayerMatchScoreTypeScore,
		MatchTeamType:    base_packet_structures.MultiplayerMatchTypeHeadToHead,
	}, sender, false)

	sender.AssignMultiplayerLobby(newLobby)

	newLobby.MultiChannel.Join(sender)

	return []string{}
}

func MpCommandInvite(sender LobbyClient, args []string) []string {
	return []string{}
}

func MpCommandLock(sender LobbyClient, args []string) []string {
	return []string{}
}

func MpCommandUnlock(sender LobbyClient, args []string) []string {
	return []string{}
}

func MpCommandSize(sender LobbyClient, args []string) []string {
	return []string{}
}

func MpCommandSet(sender LobbyClient, args []string) []string {
	return []string{}
}

func MpCommandMove(sender LobbyClient, args []string) []string {
	return []string{}
}

func MpCommandTeam(sender LobbyClient, args []string) []string {
	return []string{}
}

func MpCommandHost(sender LobbyClient, args []string) []string {
	return []string{}
}

func MpCommandSettings(sender LobbyClient, args []string) []string {
	return []string{}
}

func MpCommandStart(sender LobbyClient, args []string) []string {
	return []string{}
}

func MpCommandAbort(sender LobbyClient, args []string) []string {
	return []string{}
}

func MpCommandMap(sender LobbyClient, args []string) []string {
	return []string{}
}

func MpCommandMods(sender LobbyClient, args []string) []string {
	return []string{}
}

func MpCommandTimer(sender LobbyClient, args []string) []string {
	return []string{}
}

func MpCommandAbortTimer(sender LobbyClient, args []string) []string {
	return []string{}
}

func MpCommandKick(sender LobbyClient, args []string) []string {
	return []string{}
}

func MpCommandPassword(sender LobbyClient, args []string) []string {
	return []string{}
}

func MpCommandClose(sender LobbyClient, args []string) []string {
	return []string{}
}