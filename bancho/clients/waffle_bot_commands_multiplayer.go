package clients

import (
	"Waffle/bancho/client_manager"
	"Waffle/bancho/lobby"
	"Waffle/bancho/osu/base_packet_structures"
	"strings"
)

// WaffleBotCommandHelp !mp
func WaffleBotCommandMultiplayer(sender client_manager.WaffleClient, args []string) []string {
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

func MpCommandMake(sender client_manager.WaffleClient, args []string) []string {
	if len(args) < 2 {
		return []string{
			"!mp make: Lobby name required!",
		}
	}

	lobbyName := ""

	for i := 1; i != len(args); i++ {
		lobbyName += args[i]
	}

	newLobby := lobby.CreateNewMultiMatch(base_packet_structures.MultiplayerMatch{
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

	sender.set
}

func MpCommandInvite(sender client_manager.WaffleClient, args []string) []string {

}

func MpCommandLock(sender client_manager.WaffleClient, args []string) []string {

}

func MpCommandUnlock(sender client_manager.WaffleClient, args []string) []string {

}

func MpCommandSize(sender client_manager.WaffleClient, args []string) []string {

}

func MpCommandSet(sender client_manager.WaffleClient, args []string) []string {

}

func MpCommandMove(sender client_manager.WaffleClient, args []string) []string {

}

func MpCommandTeam(sender client_manager.WaffleClient, args []string) []string {

}

func MpCommandHost(sender client_manager.WaffleClient, args []string) []string {

}

func MpCommandSettings(sender client_manager.WaffleClient, args []string) []string {

}

func MpCommandStart(sender client_manager.WaffleClient, args []string) []string {

}

func MpCommandAbort(sender client_manager.WaffleClient, args []string) []string {

}

func MpCommandMap(sender client_manager.WaffleClient, args []string) []string {

}

func MpCommandMods(sender client_manager.WaffleClient, args []string) []string {

}

func MpCommandTimer(sender client_manager.WaffleClient, args []string) []string {

}

func MpCommandAbortTimer(sender client_manager.WaffleClient, args []string) []string {

}

func MpCommandKick(sender client_manager.WaffleClient, args []string) []string {

}

func MpCommandPassword(sender client_manager.WaffleClient, args []string) []string {

}

func MpCommandClose(sender client_manager.WaffleClient, args []string) []string {

}
