package lobby

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/database"
	"Waffle/helpers"
	"fmt"
	"sync"
)

type MultiplayerLobby struct {
	MatchId             string
	MultiChannel        *chat.Channel
	MatchInformation    base_packet_structures.MultiplayerMatch
	MatchHost           LobbyClient
	MultiClients        [8]LobbyClient
	PlayersLoaded       [8]bool
	PlayerSkipRequested [8]bool
	PlayerCompleted     [8]bool
	PlayerFailed        [8]bool
	LastScoreFrames     [8]base_packet_structures.ScoreFrame
	MatchInfoMutex      sync.Mutex
	InProgress          bool

	IrcReffed bool
	Locked    bool
}

func (multiLobby *MultiplayerLobby) LogEvent(eventType database.MatchHistoryEventType, initiator LobbyClient, extraInfo string) {
	id := int32(0)

	if initiator != nil {
		id = initiator.GetUserId()
	}

	database.LogMatchHistory(database.MatchHistoryElement{
		MatchId:        multiLobby.MatchId,
		EventType:      eventType,
		EventInitiator: uint64(id),
		ExtraInfo:      extraInfo,
	})
}

// Join gets called when a client is attempting to join the lobby
func (multiLobby *MultiplayerLobby) Join(client LobbyClient, password string) bool {
	if multiLobby.Locked {
		client.SendChatMessage("WaffleBot", "Multiplayer Lobby currently locked!", "WaffleBot")
		return false
	}

	//if they input the wrong password, join failed
	if multiLobby.MatchInformation.GamePassword != password {
		return false
	}

	multiLobby.MatchInfoMutex.Lock()

	//Inform everyone of the client, just in case they don't know them yet
	for n := 0; n != 8; n++ {
		if multiLobby.MultiClients[n] != nil {
			multiLobby.MultiClients[n].BanchoOsuUpdate(client.GetRelevantUserStats(), client.GetUserStatus())
		}
	}

	//Search for an Empty spot
	for i := 0; i != 8; i++ {
		if multiLobby.MatchInformation.SlotStatus[i] == base_packet_structures.MultiplayerMatchSlotStatusOpen {
			//Set the slot to them as well as join #multiplayer
			multiLobby.SetSlot(int32(i), client)
			multiLobby.MultiChannel.Join(client)

			multiLobby.MatchInfoMutex.Unlock()

			//Update everyone
			multiLobby.UpdateMatch()

			multiLobby.LogEvent(database.MatchHistoryEventTypeJoin, client, "")

			//Join success
			return true
		}
	}

	multiLobby.MatchInfoMutex.Unlock()

	return false
}

// SetSlot is used to set a slot to a player
func (multiLobby *MultiplayerLobby) SetSlot(slot int32, client LobbyClient) {
	//Handle for if a player is passed here, it can also be null which just sets the slot to be empty
	if client != nil {
		//Set slot nformation
		multiLobby.MatchInformation.SlotUserId[slot] = client.GetUserId()
		multiLobby.MatchInformation.SlotStatus[slot] = base_packet_structures.MultiplayerMatchSlotStatusNotReady
		multiLobby.MultiClients[slot] = client

		//Set teams, if necessary
		if multiLobby.MatchInformation.MatchTeamType == base_packet_structures.MultiplayerMatchTypeTagTeamVs || multiLobby.MatchInformation.MatchTeamType == base_packet_structures.MultiplayerMatchTypeTeamVs {
			if slot%2 == 0 {
				multiLobby.MatchInformation.SlotTeam[slot] = base_packet_structures.MultiplayerSlotTeamRed
			} else {
				multiLobby.MatchInformation.SlotTeam[slot] = base_packet_structures.MultiplayerSlotTeamBlue
			}
		}
	} else {
		//Set the slot to empty
		multiLobby.MatchInformation.SlotUserId[slot] = -1

		//If it's not locked, make it open
		if multiLobby.MatchInformation.SlotStatus[slot] != base_packet_structures.MultiplayerMatchSlotStatusLocked {
			multiLobby.MatchInformation.SlotStatus[slot] = base_packet_structures.MultiplayerMatchSlotStatusOpen
		}

		//Set team to neutral and make there be no client in that spot
		multiLobby.MatchInformation.SlotTeam[slot] = base_packet_structures.MultiplayerSlotTeamNeutral
		multiLobby.MultiClients[slot] = nil
	}
}

// MoveSlot moves a player from one slot to the other
func (multiLobby *MultiplayerLobby) MoveSlot(oldSlot int, newSlot int) {
	if oldSlot == newSlot {
		return
	}

	currentStatus := multiLobby.MatchInformation.SlotStatus[oldSlot]

	multiLobby.SetSlot(int32(newSlot), multiLobby.MultiClients[oldSlot])
	multiLobby.SetSlot(int32(oldSlot), nil)

	multiLobby.MatchInformation.SlotStatus[newSlot] = currentStatus
}

// UpdateMatch tells everyone inside the match and the lobby about the new happenings of the match
func (multiLobby *MultiplayerLobby) UpdateMatch() {
	for i := 0; i != 8; i++ {
		if multiLobby.MultiClients[i] != nil {
			multiLobby.MultiClients[i].BanchoMatchUpdate(multiLobby.MatchInformation)
		}
	}

	//Distribute in multiLobby as well
	BroadcastToLobby(func(client LobbyClient) {
		client.BanchoMatchUpdate(multiLobby.MatchInformation)
	})
}

// Handles when the IRC Referee leaves the match.
func (multiLobby *MultiplayerLobby) IrcRefereePart(client LobbyClient) {
	multiLobby.MatchInfoMutex.Lock()

	//Ignore if not IRC Reffed
	if !multiLobby.IrcReffed {
		multiLobby.MatchInfoMutex.Unlock()

		return
	}

	//Transfer host, if lobby isn't empty
	//If it's empty there's nobody to give host to
	if multiLobby.GetUsedUpSlots() != 0 {
		multiLobby.HandleHostLeave()
	}

	//If there's nobody in the multi channel, disband.
	if len(multiLobby.MultiChannel.Clients) == 0 {
		multiLobby.Disband()
	}

	//Tell everyone about it
	multiLobby.UpdateMatch()
	multiLobby.MatchInfoMutex.Unlock()
	multiLobby.LogEvent(database.MatchHistoryEventTypeLeave, client, "IRC Referee left.")
}

// Part handles a player leaving the match
func (multiLobby *MultiplayerLobby) Part(client LobbyClient) {
	multiLobby.MatchInfoMutex.Lock()

	slot := multiLobby.GetSlotFromUserId(client.GetUserId())

	//If they somehow don't exist, ignore
	if slot == -1 {
		multiLobby.MatchInfoMutex.Unlock()

		return
	}

	//Reset their slot
	multiLobby.SetSlot(int32(slot), nil)

	//If they were the host, handle that separately, as we need to pass on the host
	if multiLobby.MatchHost == client {
		multiLobby.HandleHostLeave()
	}

	//Make them leave #multiplayer
	client.BanchoChannelRevoked("#multiplayer")

	multiLobby.MultiChannel.Leave(client)

	//On IRC Reffed matches: disband when there's nobody in #multiplayer
	//On regular osu! Matches: disband when nobody's inside in the match
	if multiLobby.IrcReffed {
		if len(multiLobby.MultiChannel.Clients) == 0 {
			multiLobby.Disband()
		}
	} else {
		if multiLobby.GetUsedUpSlots() == 0 {
			multiLobby.Disband()
		}
	}

	//Tell everyone about it
	multiLobby.UpdateMatch()
	multiLobby.MatchInfoMutex.Unlock()
	multiLobby.LogEvent(database.MatchHistoryEventTypeLeave, client, "")
}

// Disband is called when everyone leaves the match
func (multiLobby *MultiplayerLobby) Disband() {
	RemoveMultiMatch(multiLobby.MatchInformation.MatchId)

	multiLobby.LogEvent(database.MatchHistoryEventTypeMatchDisbanded, nil, "")
}

// HandleHostLeave handles the host leaving, as we need to pass on the host
func (multiLobby *MultiplayerLobby) HandleHostLeave() {
	//If nobody's there anymore, disband
	if multiLobby.GetUsedUpSlots() == 0 {
		multiLobby.Disband()
	}

	//Search for a new host
	for i := 0; i != 8; i++ {
		if multiLobby.MultiClients[i] != nil {
			//If a client is found, set them to be the new host
			multiLobby.MatchHost = multiLobby.MultiClients[i]
			//Tell the new client they're host now
			multiLobby.MatchHost.BanchoMatchTransferHost()

			multiLobby.LogEvent(database.MatchHistoryEventTypeHostChange, multiLobby.MatchHost, "Host left.")

			multiLobby.MatchInformation.HostId = multiLobby.MatchHost.GetUserId()
		}
	}

	multiLobby.UpdateMatch()
}

// TryChangeSlot gets called when a player tries to change slot
func (multiLobby *MultiplayerLobby) TryChangeSlot(client LobbyClient, slotId int) {
	if multiLobby.Locked {
		client.SendChatMessage("WaffleBot", "Multiplayer Lobby currently locked! Cannot move during lock.", "WaffleBot")
		return
	}

	multiLobby.MatchInfoMutex.Lock()

	//Refuse if the slot is occupied or locked
	if multiLobby.MatchInformation.SlotStatus[slotId] == base_packet_structures.MultiplayerMatchSlotStatusLocked || (multiLobby.MatchInformation.SlotStatus[slotId]&base_packet_structures.MultiplayerMatchSlotStatusHasPlayer) > 0 {
		multiLobby.MatchInfoMutex.Unlock()

		return
	}

	//Move them to that slot and tell everyone
	multiLobby.MoveSlot(multiLobby.GetSlotFromUserId(client.GetUserId()), slotId)
	multiLobby.UpdateMatch()

	multiLobby.LogEvent(database.MatchHistoryEventTypeMove, client, fmt.Sprintf("Moved to slot %d", slotId))

	multiLobby.MatchInfoMutex.Unlock()
}

// ChangeTeam gets called when a player is trying to change their team
func (multiLobby *MultiplayerLobby) ChangeTeam(client LobbyClient) {
	if multiLobby.Locked {
		client.SendChatMessage("WaffleBot", "Multiplayer Lobby currently locked! Cannot change teams during lock.", "WaffleBot")
		return
	}

	multiLobby.MatchInfoMutex.Lock()

	clientSlot := multiLobby.GetSlotFromUserId(client.GetUserId())

	if clientSlot == -1 {
		multiLobby.MatchInfoMutex.Unlock()

		return
	}

	color := "Red"

	//Flip colors
	if multiLobby.MatchInformation.SlotTeam[clientSlot] == base_packet_structures.MultiplayerSlotTeamRed {
		multiLobby.MatchInformation.SlotTeam[clientSlot] = base_packet_structures.MultiplayerSlotTeamBlue

		color = "Blue"
	} else if multiLobby.MatchInformation.SlotTeam[clientSlot] == base_packet_structures.MultiplayerSlotTeamBlue {
		multiLobby.MatchInformation.SlotTeam[clientSlot] = base_packet_structures.MultiplayerSlotTeamRed
	} else {
		multiLobby.MatchInformation.SlotTeam[clientSlot] = base_packet_structures.MultiplayerSlotTeamRed
	}

	//Tell everyone
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()

	multiLobby.LogEvent(database.MatchHistoryEventTypeChangeTeam, client, fmt.Sprintf("Is now Team %s", color))
}

// TransferHost gets called when the host willingly gives up their host
func (multiLobby *MultiplayerLobby) TransferHost(client LobbyClient, slotId int) {
	multiLobby.MatchInfoMutex.Lock()

	if multiLobby.MatchHost != client {
		multiLobby.MatchInfoMutex.Unlock()

		return
	}

	if multiLobby.Locked {
		client.SendChatMessage("WaffleBot", "Multiplayer Lobby currently locked! Cannot transfer host during lock.", "WaffleBot")

		multiLobby.MatchInfoMutex.Unlock()
		return
	}

	//set the new host
	multiLobby.MatchHost = multiLobby.MultiClients[slotId]
	multiLobby.MatchInformation.HostId = multiLobby.MatchHost.GetUserId()

	//Tell them about it
	multiLobby.MatchHost.BanchoMatchTransferHost()

	//Tell everyone
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()

	multiLobby.LogEvent(database.MatchHistoryEventTypeHostChange, client, fmt.Sprintf("New host is now UserID: %d; Username %s", multiLobby.MatchHost.GetUserId(), multiLobby.MatchHost.GetUsername()))
}

// ReadyUp gets called when a player has clicked the Ready button
func (multiLobby *MultiplayerLobby) ReadyUp(client LobbyClient) {
	if multiLobby.Locked {
		client.SendChatMessage("WaffleBot", "Multiplayer Lobby currently locked! Cannot change ready state during lock.", "WaffleBot")
		return
	}

	multiLobby.MatchInfoMutex.Lock()

	clientSlot := multiLobby.GetSlotFromUserId(client.GetUserId())

	if clientSlot == -1 {
		multiLobby.MatchInfoMutex.Unlock()

		return
	}

	//Set them to be ready and tell everyone they're ready
	multiLobby.MatchInformation.SlotStatus[clientSlot] = base_packet_structures.MultiplayerMatchSlotStatusReady
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()

	multiLobby.LogEvent(database.MatchHistoryEventTypeReady, client, "")
}

// Unready gets called when a player has changed their mind about being ready and pressed the not ready button
func (multiLobby *MultiplayerLobby) Unready(client LobbyClient) {
	if multiLobby.Locked {
		client.SendChatMessage("WaffleBot", "Multiplayer Lobby currently locked! Cannot change ready state during lock.", "WaffleBot")
		return
	}

	multiLobby.MatchInfoMutex.Lock()

	clientSlot := multiLobby.GetSlotFromUserId(client.GetUserId())

	if clientSlot == -1 {
		multiLobby.MatchInfoMutex.Unlock()

		return
	}

	//Set them to be not ready and tell everyone
	multiLobby.MatchInformation.SlotStatus[clientSlot] = base_packet_structures.MultiplayerMatchSlotStatusNotReady
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()

	multiLobby.LogEvent(database.MatchHistoryEventTypeUnready, client, "")
}

// ChangeSettings gets called when the host of the lobby changes some settings
func (multiLobby *MultiplayerLobby) ChangeSettings(client LobbyClient, matchSettings base_packet_structures.MultiplayerMatch) {
	multiLobby.MatchInfoMutex.Lock()

	if multiLobby.MatchHost != client {
		multiLobby.MatchInfoMutex.Unlock()

		return
	}

	if multiLobby.Locked {
		client.SendChatMessage("WaffleBot", "Multiplayer Lobby currently locked! Cannot change ready state during lock.", "WaffleBot")

		multiLobby.MatchInfoMutex.Unlock()
		return
	}

	//We're building a diff for the match logging.
	diff := ""

	if multiLobby.MatchInformation.ActiveMods != matchSettings.ActiveMods {
		oldMods := helpers.FormatMods(uint32(multiLobby.MatchInformation.ActiveMods))
		newMods := helpers.FormatMods(uint32(matchSettings.ActiveMods))

		diff += fmt.Sprintf("Mods changed (was: %s; is: %s)\n", oldMods, newMods)
	}

	if multiLobby.MatchInformation.GameName != matchSettings.GameName {
		diff += fmt.Sprintf("Match name changed (was: %s; is %s)\n", multiLobby.MatchInformation.GameName, matchSettings.GameName)
	}

	if multiLobby.MatchInformation.GamePassword != matchSettings.GamePassword {
		diff += "Password was changed.\n"
	}

	if multiLobby.MatchInformation.BeatmapChecksum != matchSettings.BeatmapChecksum {
		diff += fmt.Sprintf("Map was changed to Beatmap ID: %d; Name: %s\n", matchSettings.BeatmapId, matchSettings.BeatmapName)
	}

	if multiLobby.MatchInformation.Playmode != matchSettings.Playmode {
		oldMode := helpers.FormatPlaymodes(multiLobby.MatchInformation.Playmode)
		newMode := helpers.FormatPlaymodes(matchSettings.Playmode)

		diff += fmt.Sprintf("Playmode changed (was: %s; is %s)\n", oldMode, newMode)
	}

	if multiLobby.MatchInformation.MatchScoringType != matchSettings.MatchScoringType {
		oldMode := helpers.FormatScoringType(multiLobby.MatchInformation.MatchScoringType)
		newMode := helpers.FormatScoringType(matchSettings.MatchScoringType)

		diff += fmt.Sprintf("Scoring type changed (was: %s; is %s)\n", oldMode, newMode)
	}

	if multiLobby.MatchInformation.MatchTeamType != matchSettings.MatchTeamType {
		oldMode := helpers.FormatMatchTeamTypes(multiLobby.MatchInformation.MatchTeamType)
		newMode := helpers.FormatMatchTeamTypes(matchSettings.MatchTeamType)

		diff += fmt.Sprintf("Team type changed (was: %s; is %s)\n", oldMode, newMode)
	}

	//Update the settings and tell everyone
	multiLobby.MatchInformation = matchSettings
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()

	multiLobby.LogEvent(database.MatchHistoryEventTypeSettingsChanged, multiLobby.MatchHost, diff)
}

// ChangeMods gets called when the host of the lobby changes which mods are going to get played
func (multiLobby *MultiplayerLobby) ChangeMods(client LobbyClient, newMods int32) {
	multiLobby.MatchInfoMutex.Lock()

	if multiLobby.MatchHost != client {
		multiLobby.MatchInfoMutex.Unlock()

		return
	}

	if multiLobby.Locked {
		client.SendChatMessage("WaffleBot", "Multiplayer Lobby currently locked! Cannot change ready state during lock.", "WaffleBot")

		multiLobby.MatchInfoMutex.Unlock()
		return
	}

	oldMods := helpers.FormatMods(uint32(multiLobby.MatchInformation.ActiveMods))

	//Set new mods and tell everyone
	multiLobby.MatchInformation.ActiveMods = uint16(newMods)
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()

	newModsFmt := helpers.FormatMods(uint32(newMods))

	multiLobby.LogEvent(database.MatchHistoryEventTypeModsChanged, multiLobby.MatchHost, fmt.Sprintf("was %s; is %s", oldMods, newModsFmt))
}

// LockSlot gets called when the host attempts to lock/unlock a slot
func (multiLobby *MultiplayerLobby) LockSlot(client LobbyClient, slotId int) {
	multiLobby.MatchInfoMutex.Lock()

	if multiLobby.MatchHost != client {
		multiLobby.MatchInfoMutex.Unlock()

		return
	}

	if multiLobby.Locked {
		client.SendChatMessage("WaffleBot", "Multiplayer Lobby currently locked! Cannot lock slots during lock.", "WaffleBot")

		multiLobby.MatchInfoMutex.Unlock()
		return
	}

	//don't allow the host to kick themselves by locking their slot
	if multiLobby.MultiClients[slotId] == multiLobby.MatchHost {
		multiLobby.MatchInfoMutex.Unlock()

		return
	}

	//If we lock a slot with a player inside, we kick them
	if (multiLobby.MatchInformation.SlotStatus[slotId] & base_packet_structures.MultiplayerMatchSlotStatusHasPlayer) > 0 {
		droppedClient := multiLobby.MultiClients[slotId]

		multiLobby.MatchInfoMutex.Unlock()

		droppedClient.LeaveCurrentMatch()
		droppedClient.BanchoMatchUpdate(multiLobby.MatchInformation)

		multiLobby.MatchInfoMutex.Lock()

		multiLobby.LogEvent(database.MatchHistoryEventTypeKick, multiLobby.MatchHost, "Slot was locked")
	}

	//If it's locked already, make it open
	if multiLobby.MatchInformation.SlotStatus[slotId] == base_packet_structures.MultiplayerMatchSlotStatusLocked {
		multiLobby.MatchInformation.SlotStatus[slotId] = base_packet_structures.MultiplayerMatchSlotStatusOpen

		multiLobby.UpdateMatch()
		multiLobby.MatchInfoMutex.Unlock()

		multiLobby.LogEvent(database.MatchHistoryEventTypeUnlock, multiLobby.MatchHost, fmt.Sprintf("Slot %d was unlocked", slotId))

		return
	}

	//Don't allow all slots to be locked
	if multiLobby.GetOpenSlotCount() > 2 && multiLobby.MatchInformation.SlotStatus[slotId] == base_packet_structures.MultiplayerMatchSlotStatusOpen {
		multiLobby.MatchInformation.SlotStatus[slotId] = base_packet_structures.MultiplayerMatchSlotStatusLocked

		multiLobby.LogEvent(database.MatchHistoryEventTypeLock, multiLobby.MatchHost, fmt.Sprintf("Slot %d was locked", slotId))
	}

	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()

}

// InformNoBeatmap gets called when a player happens to be missing the map thats about to be played
func (multiLobby *MultiplayerLobby) InformNoBeatmap(client LobbyClient) {
	multiLobby.MatchInfoMutex.Lock()

	slot := multiLobby.GetSlotFromUserId(client.GetUserId())

	if slot == -1 {
		multiLobby.MatchInfoMutex.Unlock()

		return
	}

	//Mark them as missing the map and tell everyone
	multiLobby.MatchInformation.SlotStatus[slot] = base_packet_structures.MultiplayerMatchSlotStatusMissingMap
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()
}

// InformGotBeatmap gets called whenever a player has now gotten the beatmap that they were missing earlier
func (multiLobby *MultiplayerLobby) InformGotBeatmap(client LobbyClient) {
	multiLobby.MatchInfoMutex.Lock()

	slot := multiLobby.GetSlotFromUserId(client.GetUserId())

	if slot == -1 {
		multiLobby.MatchInfoMutex.Unlock()

		return
	}

	//Set them to be not ready and tell everyone
	multiLobby.MatchInformation.SlotStatus[slot] = base_packet_structures.MultiplayerMatchSlotStatusNotReady
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()
}

// InformLoadComplete gets called when a player has loaded into the game
func (multiLobby *MultiplayerLobby) InformLoadComplete(client LobbyClient) {
	multiLobby.MatchInfoMutex.Lock()

	slot := multiLobby.GetSlotFromUserId(client.GetUserId())

	if slot == -1 {
		multiLobby.MatchInfoMutex.Unlock()

		return
	}

	//Set their slot to be fully loaded
	multiLobby.PlayersLoaded[slot] = true

	//Check if everyone has loaded in, if yes then tell everyone that everyone's ready and begin!
	if multiLobby.HaveAllPlayersLoaded() {
		for i := 0; i != 8; i++ {
			if multiLobby.MultiClients[i] != nil {
				multiLobby.MultiClients[i].BanchoMatchAllPlayersLoaded()
			}
		}

		multiLobby.LogEvent(database.MatchHistoryEventTypePlayingStarted, nil, "")
	}

	multiLobby.MatchInfoMutex.Unlock()
}

// InformScoreUpdate this happens every time a player hits a circle or gets a slidertick or whatever
func (multiLobby *MultiplayerLobby) InformScoreUpdate(client LobbyClient, scoreFrame base_packet_structures.ScoreFrame) {
	multiLobby.MatchInfoMutex.Lock()

	slot := multiLobby.GetSlotFromUserId(client.GetUserId())

	if slot == -1 {
		multiLobby.MatchInfoMutex.Unlock()

		return
	}

	//Set their slot id
	scoreFrame.Id = uint8(slot)
	//Currently unused, but could be useful to display statistics after the match had ended and stuff
	multiLobby.LastScoreFrames[slot] = scoreFrame

	//Tell everyone about their new score
	for i := 0; i != 8; i++ {
		if multiLobby.MultiClients[i] != nil {
			multiLobby.MultiClients[i].BanchoMatchScoreUpdate(scoreFrame)
		}
	}

	multiLobby.MatchInfoMutex.Unlock()
}

// InformCompletion gets called whenever a client has finished playing a map
func (multiLobby *MultiplayerLobby) InformCompletion(client LobbyClient) {
	multiLobby.MatchInfoMutex.Lock()

	slot := multiLobby.GetSlotFromUserId(client.GetUserId())

	if slot == -1 {
		multiLobby.MatchInfoMutex.Unlock()

		return
	}

	//Set them to be completed
	multiLobby.PlayerCompleted[slot] = true

	//Check if everyone completed
	if multiLobby.HaveAllPlayersCompleted() {
		//Set the match to no longer be in progress
		multiLobby.InProgress = false

		for i := 0; i != 8; i++ {
			//Reset all states
			multiLobby.PlayerCompleted[i] = false
			multiLobby.PlayerSkipRequested[i] = false
			multiLobby.PlayersLoaded[i] = false
			multiLobby.PlayerFailed[i] = false

			if multiLobby.MultiClients[i] != nil {
				multiLobby.MatchInformation.SlotStatus[i] = base_packet_structures.MultiplayerMatchSlotStatusNotReady

				multiLobby.MultiClients[i].BanchoMatchComplete()
			}
		}
	}

	//Tell everyone
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()

	info := fmt.Sprintf("Score: %d\n", multiLobby.LastScoreFrames[slot].TotalScore)
	info += fmt.Sprintf("Max Combo: %d\n", multiLobby.LastScoreFrames[slot].MaxCombo)
	info += fmt.Sprintf("Final Combo: %d\n", multiLobby.LastScoreFrames[slot].CurrentCombo)
	info += fmt.Sprintf("Perfect: %t\n", multiLobby.LastScoreFrames[slot].Perfect)
	info += fmt.Sprintf("300: %d\n", multiLobby.LastScoreFrames[slot].Count300)
	info += fmt.Sprintf("100: %d\n", multiLobby.LastScoreFrames[slot].Count100)
	info += fmt.Sprintf("50: %d\n", multiLobby.LastScoreFrames[slot].Count50)
	info += fmt.Sprintf("Miss: %d\n", multiLobby.LastScoreFrames[slot].CountMiss)
	info += fmt.Sprintf("Geki: %d\n", multiLobby.LastScoreFrames[slot].CountGeki)
	info += fmt.Sprintf("Katu: %d", multiLobby.LastScoreFrames[slot].CountKatu)

	multiLobby.LogEvent(database.MatchHistoryEventTypeFinalScore, client, info)
}

// InformPressedSkip gets called when a player pressed skip in multi
func (multiLobby *MultiplayerLobby) InformPressedSkip(client LobbyClient) {
	multiLobby.MatchInfoMutex.Lock()

	slot := multiLobby.GetSlotFromUserId(client.GetUserId())

	if slot == -1 {
		multiLobby.MatchInfoMutex.Unlock()

		return
	}

	//Set their slot to be skipped
	multiLobby.PlayerSkipRequested[slot] = true

	//Tell everyone that they skipped
	for i := 0; i != 8; i++ {
		if multiLobby.MultiClients[i] != nil {
			multiLobby.MultiClients[i].BanchoMatchPlayerSkipped(int32(slot))
		}
	}

	//If everyone skipped, tell everyone that it's okay to skip
	if multiLobby.HaveAllPlayersSkipped() {
		for i := 0; i != 8; i++ {
			if multiLobby.MultiClients[i] != nil {
				multiLobby.MultiClients[i].BanchoMatchSkip()
			}
		}
	}

	multiLobby.MatchInfoMutex.Unlock()
}

// InformFailed gets called whenever a client fails
func (multiLobby *MultiplayerLobby) InformFailed(client LobbyClient) {
	multiLobby.MatchInfoMutex.Lock()

	slot := multiLobby.GetSlotFromUserId(client.GetUserId())

	if slot == -1 {
		multiLobby.MatchInfoMutex.Unlock()

		return
	}

	//Set them as failed
	multiLobby.PlayerFailed[slot] = true

	//Tell everyone they failed
	for i := 0; i != 8; i++ {
		if multiLobby.MultiClients[i] != nil {
			multiLobby.MultiClients[i].BanchoMatchPlayerFailed(int32(slot))
		}
	}

	multiLobby.MatchInfoMutex.Unlock()

	multiLobby.LogEvent(database.MatchHistoryEventTypePlayerFail, client, "")
}

// StartGame gets called whenever the host starts the game
func (multiLobby *MultiplayerLobby) StartGame(client LobbyClient) {
	multiLobby.MatchInfoMutex.Lock()

	if multiLobby.MatchHost != client {
		multiLobby.MatchInfoMutex.Unlock()

		return
	}

	//Sets the game to be in progress
	multiLobby.InProgress = true

	//Tell everyone to start
	for i := 0; i != 8; i++ {
		if multiLobby.MultiClients[i] != nil {
			multiLobby.MatchInformation.SlotStatus[i] = base_packet_structures.MultiplayerMatchSlotStatusPlaying

			multiLobby.MultiClients[i].BanchoMatchStart(multiLobby.MatchInformation)
		}
	}

	//Tell everyone, in lobby aswell
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()

	multiLobby.LogEvent(database.MatchHistoryEventTypeMatchStarted, client, "")
}

// Various IRC Ref related things:

func (multiLobby *MultiplayerLobby) RefereeLock(client LobbyClient) {
	if client != multiLobby.MatchHost {
		return
	}

	multiLobby.Locked = true

	multiLobby.LogEvent(database.MatchHistoryEventTypeMatchRefLocked, client, "Referee locked the lobby.")
}

func (multiLobby *MultiplayerLobby) RefereeUnlock(client LobbyClient) {
	if client != multiLobby.MatchHost {
		return
	}

	multiLobby.Locked = false

	multiLobby.LogEvent(database.MatchHistoryEventTypeMatchRefUnlocked, client, "Referee unlocked the lobby.")
}

func (multiLobby *MultiplayerLobby) MovePlayerUp(client LobbyClient, slotToMove int) bool {
	if slotToMove == 0 {
		return false
	}

	if multiLobby.MatchInformation.SlotStatus[slotToMove] == base_packet_structures.MultiplayerMatchSlotStatusOpen {
		return false
	}

	attemptedSlot := multiLobby.MatchInformation.SlotStatus[slotToMove-1]

	if attemptedSlot == base_packet_structures.MultiplayerMatchSlotStatusOpen {
		multiLobby.MoveSlot(slotToMove, slotToMove-1)

		return true
	}

	return false
}

func (multiLobby *MultiplayerLobby) IsTight() bool {
	consecutive := 0
	for i := 0; i != 8; i++ {
		if multiLobby.MatchInformation.SlotStatus[i] != base_packet_structures.MultiplayerMatchSlotStatusOpen {
			consecutive++
		} else {
			break
		}
	}

	return consecutive == multiLobby.GetUsedUpSlots()
}

func (multiLobby *MultiplayerLobby) SetSize(client LobbyClient, size int) {
	if client != multiLobby.MatchHost {
		return
	}

	if multiLobby.Locked {
		client.SendChatMessage("WaffleBot", "Multiplayer Lobby currently locked! Cannot change lobby arrangement during lock.", "WaffleBot")

		return
	}

	multiLobby.MatchInfoMutex.Lock()

	//Move everyone as tight as possible
	for !multiLobby.IsTight() {
		for i := 1; i != 8; i++ {
			multiLobby.MovePlayerUp(client, i)
		}
	}

	//Lock all the bottom slots
	for i := 7; i != 7-(8-size); i-- {
		multiLobby.MatchInformation.SlotStatus[i] = base_packet_structures.MultiplayerMatchSlotStatusLocked
	}

	//Tell everyone, in lobby aswell
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()
}
