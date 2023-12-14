package lobby

import (
	"Waffle/bancho/osu/base_packet_structures"
	"fmt"
	"strconv"
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
		SlotStatus:       [8]uint8{1, 1, 1, 1, 1, 1, 1, 1},
	}, sender, true)

	sender.AssignMultiplayerLobby(newLobby)

	newLobby.MultiChannel.Join(sender)
	sender.AddJoinedChannel(newLobby.MultiChannel)

	return []string{}
}

func MpCommandInvite(sender LobbyClient, args []string) []string {
	currentLobby := sender.GetMultiplayerLobby()
	if currentLobby == nil {
		return []string{
			"!mp invite: Only usable inside multiplayer lobby!",
		}
	}

	if len(args) < 2 {
		return []string{
			"!mp invite: Username required!",
		}
	}

	username := ""

	for i := 1; i != len(args); i++ {
		username += args[i]
	}

	//TODO: figure out how to do this

	return []string{}
}

func MpCommandLock(sender LobbyClient, args []string) []string {
	currentLobby := sender.GetMultiplayerLobby()
	if currentLobby == nil {
		return []string{
			"!mp lock: Only usable inside multiplayer lobby!",
		}
	}

	currentLobby.RefereeLock(sender)

	return []string{}
}

func MpCommandUnlock(sender LobbyClient, args []string) []string {
	currentLobby := sender.GetMultiplayerLobby()
	if currentLobby == nil {
		return []string{
			"!mp unlock: Only usable inside multiplayer lobby!",
		}
	}

	currentLobby.RefereeUnlock(sender)

	return []string{}
}

func MpCommandSize(sender LobbyClient, args []string) []string {
	currentLobby := sender.GetMultiplayerLobby()
	if currentLobby == nil {
		return []string{
			"!mp size: Only usable inside multiplayer lobby!",
		}
	}

	if len(args) < 2 {
		return []string{
			"!mp invite: Size required!",
		}
	}

	size := args[1]

	num, err := strconv.ParseInt(size, 10, 64)

	if err != nil {
		return []string{
			"!mp size: make sure the size is a number.",
		}
	}

	if currentLobby.GetUsedUpSlots() > int(num) {
		return []string{
			"!mp size: there are more used up slots than you want to size down to.",
		}
	}

	currentLobby.SetSize(sender, int(num))

	return []string{}
}

func MpCommandSet(sender LobbyClient, args []string) []string {
	currentLobby := sender.GetMultiplayerLobby()
	if currentLobby == nil {
		return []string{
			"!mp set: Only usable inside multiplayer lobby!",
		}
	}

	if len(args) < 2 {
		return []string{
			"!mp invite: Username required!",
		}
	}

	teamMode := args[1]
	scoreMode := ""

	if len(args) >= 3 {
		scoreMode = args[2]
	}

	newTeamType := currentLobby.MatchInformation.MatchTeamType

	switch teamMode {
	case "0":
		newTeamType = base_packet_structures.MultiplayerMatchTypeHeadToHead
	case "1":
		newTeamType = base_packet_structures.MultiplayerMatchTypeTagCoop
	case "2":
		newTeamType = base_packet_structures.MultiplayerMatchTypeTeamVs
	case "3":
		newTeamType = base_packet_structures.MultiplayerMatchTypeTagTeamVs
	}

	newScoringMode := currentLobby.MatchInformation.MatchScoringType

	if scoreMode != "" {
		switch scoreMode {
		case "0":
			newScoringMode = base_packet_structures.MultiplayerMatchScoreTypeScore
		case "1":
			newScoringMode = base_packet_structures.MultiplayerMatchScoreTypeAccuracy
		}
	}

	matchSetings := currentLobby.MatchInformation
	matchSetings.MatchTeamType = newTeamType
	matchSetings.MatchScoringType = newScoringMode

	currentLobby.ChangeSettings(sender, matchSetings)

	return []string{}
}

func MpCommandMove(sender LobbyClient, args []string) []string {
	currentLobby := sender.GetMultiplayerLobby()
	if currentLobby == nil {
		return []string{
			"!mp move: Only usable inside multiplayer lobby!",
		}
	}

	if len(args) < 3 {
		return []string{
			"!mp move: Username and slot required!",
		}
	}

	username := args[1]
	slot := args[2]

	parsedSlot, err := strconv.ParseInt(slot, 10, 64)

	if err != nil {
		return []string{
			"!mp move: Actual number for Slot required.",
		}
	}

	if parsedSlot > 7 || parsedSlot < 0 {
		return []string{
			"!mp move: Slot outside range.",
		}
	}

	for i := 0; i != 8; i++ {
		currentClient := currentLobby.MultiClients[i]

		if currentClient == nil {
			continue
		}

		if currentClient.GetUsername() == username {
			currentLobby.TryChangeSlot(currentLobby.MultiClients[i], int(parsedSlot))
		}
	}

	return []string{}
}

func MpCommandTeam(sender LobbyClient, args []string) []string {
	currentLobby := sender.GetMultiplayerLobby()
	if currentLobby == nil {
		return []string{
			"!mp team: Only usable inside multiplayer lobby!",
		}
	}

	if len(args) < 3 {
		return []string{
			"!mp move: Username and team color required!",
		}
	}

	username := args[1]
	team := args[2]

	actualTeam := base_packet_structures.MultiplayerSlotTeamRed

	switch team {
	case "red":
		actualTeam = base_packet_structures.MultiplayerSlotTeamRed
	case "blue":
		actualTeam = base_packet_structures.MultiplayerSlotTeamBlue
	default:
		return []string{
			fmt.Sprintf("!mp team: %s is not a valid team.", team),
		}
	}

	for i := 0; i != 8; i++ {
		currentClient := currentLobby.MultiClients[i]

		if currentClient == nil {
			continue
		}

		if currentClient.GetUsername() == username {
			slot := currentLobby.GetSlotFromUserId(currentClient.GetUserId())

			currentLobby.MatchInformation.SlotTeam[slot] = actualTeam
			currentLobby.UpdateMatch()
		}
	}

	return []string{}
}

func MpCommandHost(sender LobbyClient, args []string) []string {
	currentLobby := sender.GetMultiplayerLobby()
	if currentLobby == nil {
		return []string{
			"!mp host: Only usable inside multiplayer lobby!",
		}
	}

	if len(args) < 2 {
		return []string{
			"!mp host: Username required!",
		}
	}

	username := args[1]

	for i := 0; i != 8; i++ {
		currentClient := currentLobby.MultiClients[i]

		if currentClient == nil {
			continue
		}

		if currentClient.GetUsername() == username {
			slot := currentLobby.GetSlotFromUserId(currentClient.GetUserId())

			currentLobby.TransferHost(sender, slot)
			currentLobby.UpdateMatch()
		}
	}

	return []string{}
}

func MpCommandSettings(sender LobbyClient, args []string) []string {
	currentLobby := sender.GetMultiplayerLobby()
	if currentLobby == nil {
		return []string{
			"!mp settings: Only usable inside multiplayer lobby!",
		}
	}

	return []string{}
}

func MpCommandStart(sender LobbyClient, args []string) []string {
	currentLobby := sender.GetMultiplayerLobby()
	if currentLobby == nil {
		return []string{
			"!mp start: Only usable inside multiplayer lobby!",
		}
	}

	return []string{}
}

func MpCommandAbort(sender LobbyClient, args []string) []string {
	currentLobby := sender.GetMultiplayerLobby()
	if currentLobby == nil {
		return []string{
			"!mp abort: Only usable inside multiplayer lobby!",
		}
	}

	return []string{}
}

func MpCommandMap(sender LobbyClient, args []string) []string {
	currentLobby := sender.GetMultiplayerLobby()
	if currentLobby == nil {
		return []string{
			"!mp map: Only usable inside multiplayer lobby!",
		}
	}

	return []string{}
}

func MpCommandMods(sender LobbyClient, args []string) []string {
	currentLobby := sender.GetMultiplayerLobby()
	if currentLobby == nil {
		return []string{
			"!mp mods: Only usable inside multiplayer lobby!",
		}
	}

	return []string{}
}

func MpCommandTimer(sender LobbyClient, args []string) []string {
	currentLobby := sender.GetMultiplayerLobby()
	if currentLobby == nil {
		return []string{
			"!mp timer: Only usable inside multiplayer lobby!",
		}
	}

	return []string{}
}

func MpCommandAbortTimer(sender LobbyClient, args []string) []string {
	currentLobby := sender.GetMultiplayerLobby()
	if currentLobby == nil {
		return []string{
			"!mp abort: Only usable inside multiplayer lobby!",
		}
	}

	return []string{}
}

func MpCommandKick(sender LobbyClient, args []string) []string {
	currentLobby := sender.GetMultiplayerLobby()
	if currentLobby == nil {
		return []string{
			"!mp kick: Only usable inside multiplayer lobby!",
		}
	}

	return []string{}
}

func MpCommandPassword(sender LobbyClient, args []string) []string {
	currentLobby := sender.GetMultiplayerLobby()
	if currentLobby == nil {
		return []string{
			"!mp password: Only usable inside multiplayer lobby!",
		}
	}

	return []string{}
}

func MpCommandClose(sender LobbyClient, args []string) []string {
	currentLobby := sender.GetMultiplayerLobby()
	if currentLobby == nil {
		return []string{
			"!mp close: Only usable inside multiplayer lobby!",
		}
	}

	return []string{}
}
