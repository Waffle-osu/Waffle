package lobby

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/database"
	"Waffle/helpers"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
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

	messages := []string{
		fmt.Sprintf("-- Lobby name: %s", currentLobby.MatchInformation.GameName),
		fmt.Sprintf("Has password: %t", currentLobby.MatchInformation.GamePassword != ""),
		fmt.Sprintf("Beatmap Name: %s", currentLobby.MatchInformation.BeatmapName),
		fmt.Sprintf("Beatmap ID: %d", currentLobby.MatchInformation.BeatmapId),
		fmt.Sprintf("Playmode: %s", helpers.FormatPlaymodes(currentLobby.MatchInformation.Playmode)),
		fmt.Sprintf("Team Type: %s", helpers.FormatMatchTeamTypes(currentLobby.MatchInformation.MatchTeamType)),
		fmt.Sprintf("Scoring Type: %s", helpers.FormatScoringType(currentLobby.MatchInformation.MatchScoringType)),
		fmt.Sprintf("Host: %s", currentLobby.MatchHost.GetUsername()),
		"-- Slots:",
	}

	for i := 0; i != 8; i++ {
		formattedSlot := ""

		if currentLobby.MatchInformation.SlotStatus[i] == base_packet_structures.MultiplayerMatchSlotStatusLocked {
			formattedSlot = fmt.Sprintf("[%d] Locked.", i)
		} else {
			if currentLobby.MatchInformation.SlotStatus[i] == base_packet_structures.MultiplayerMatchSlotStatusOpen {
				formattedSlot = fmt.Sprintf("[%d] Open..", i)
			} else {
				userId := currentLobby.MatchInformation.SlotUserId[i]

				for j := 0; j != 8; j++ {
					currentClient := currentLobby.MultiClients[j]

					if currentClient == nil {
						continue
					}

					if currentClient.GetUserId() == userId {
						formattedSlot = fmt.Sprintf("[%d] Name: %s; Status: %s", i, currentClient.GetUsername(), helpers.FormatSlotStatus(currentLobby.MatchInformation.SlotStatus[i]))
					}
				}
			}
		}

		messages = append(messages, formattedSlot)
	}

	return messages
}

func timeTicker(countdown int, tickerMessagePrefix string, tickerMessageSender chat.ChatClient, matchHost LobbyClient, ctx context.Context, onDone func(sender LobbyClient)) {
	ticker := time.NewTicker(1 * time.Second)
	toStart := countdown

	sendMsg := func(message string) {
		channel := matchHost.GetMultiplayerLobby().MultiChannel
		channel.SendMessage(tickerMessageSender, message, "#multiplayer")
	}

	send := func(time int) {
		millis := uint64(toStart * 1000)
		timeFormatted := helpers.FormatTime(millis)

		sendMsg(fmt.Sprintf("%s in %s...", tickerMessagePrefix, timeFormatted))
	}

	for {
		select {
		case <-ticker.C:
			toStart--

			switch toStart {
			case 1800:
				send(toStart)
			case 600:
				send(toStart)
			case 300:
				send(toStart)
			case 60:
				send(toStart)
			case 30:
				send(toStart)
			case 10:
				send(toStart)
			case 5:
				send(toStart)
			case 4:
				send(toStart)
			case 3:
				send(toStart)
			case 2:
				send(toStart)
			case 1:
				send(toStart)
			case 0:
				ticker.Stop()
				onDone(matchHost)
				return
			}
		case <-ctx.Done():
			ticker.Stop()

			sendMsg("Timer aborted!")

			return
		}
	}
}

func MpCommandStart(sender LobbyClient, args []string) []string {
	currentLobby := sender.GetMultiplayerLobby()
	if currentLobby == nil {
		return []string{
			"!mp start: Only usable inside multiplayer lobby!",
		}
	}

	startTime := 0

	if len(args) >= 2 {
		parsedTime, err := strconv.ParseInt(args[1], 10, 64)

		if err != nil {
			return []string{
				"!mp start: Start time must be a number!",
			}
		}

		startTime = int(parsedTime)
	}

	if startTime > 0 {
		ctx, cancel := context.WithCancel(context.Background())

		currentLobby.MatchStartCancel = cancel

		onDone := func(sender LobbyClient) {
			sender.GetMultiplayerLobby().StartGame(sender)
		}

		go timeTicker(startTime, "Starting game", LobbyWaffleBot{}, sender, ctx, onDone)

		return []string{
			fmt.Sprintf("Starting game in %d seconds", startTime),
		}
	} else {
		//Just start the game
		currentLobby.StartGame(sender)
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

	if currentLobby.MatchStartCancel != nil {
		currentLobby.MatchStartCancel()
	}

	return []string{
		"Match start countdown cancelled.",
	}
}

func MpCommandMap(sender LobbyClient, args []string) []string {
	currentLobby := sender.GetMultiplayerLobby()
	if currentLobby == nil {
		return []string{
			"!mp map: Only usable inside multiplayer lobby!",
		}
	}

	if len(args) < 2 {
		return []string{
			"!mp map: Beatmap ID required!",
		}
	}

	id, err := strconv.ParseInt(args[1], 10, 64)

	if err != nil {
		return []string{
			"!mp map: Beatmap ID must be a number",
		}
	}

	queryResult, beatmap := database.BeatmapsGetById(int32(id))
	queryResultSet, beatmapset := database.BeatmapsetsGetBeatmapsetById(beatmap.BeatmapsetId)

	if queryResult != 0 || queryResultSet != 0 {
		return []string{
			"!mp map: Beatmap Query failed!",
		}
	}

	newSettings := currentLobby.MatchInformation

	newSettings.BeatmapId = beatmap.BeatmapId
	newSettings.BeatmapChecksum = beatmap.BeatmapMd5
	newSettings.BeatmapName = fmt.Sprintf("%s - %s [%s]", beatmapset.Artist, beatmapset.Title, beatmap.Version)

	currentLobby.ChangeSettings(sender, newSettings)

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
			"!mp start: Only usable inside multiplayer lobby!",
		}
	}

	startTime := 30

	if len(args) >= 2 {
		parsedTime, err := strconv.ParseInt(args[1], 10, 64)

		if err != nil {
			return []string{
				"!mp start: Start time must be a number!",
			}
		}

		startTime = int(parsedTime)
	}

	if startTime > 0 {
		ctx, cancel := context.WithCancel(context.Background())

		currentLobby.TimerCancel = cancel

		onDone := func(sender LobbyClient) {
			currentLobby.MultiChannel.SendMessage(LobbyWaffleBot{}, "Countdown finished.", "#multiplayer")
		}

		go timeTicker(startTime, "Countdown ends", LobbyWaffleBot{}, sender, ctx, onDone)

		return []string{
			fmt.Sprintf("Started countdown for %d seconds", startTime),
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

	if currentLobby.TimerCancel != nil {
		currentLobby.TimerCancel()
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
