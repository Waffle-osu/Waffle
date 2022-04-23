package lobby

import (
	"Waffle/waffle/chat"
	"Waffle/waffle/packets"
	"sync"
)

type MultiplayerLobby struct {
	MultiChannel        *chat.Channel
	MatchInformation    packets.MultiplayerMatch
	MatchHost           LobbyClient
	MultiClients        [8]LobbyClient
	PlayersLoaded       [8]bool
	PlayerSkipRequested [8]bool
	PlayerCompleted     [8]bool
	LastScoreFrames     [8]packets.ScoreFrame
	MatchInfoMutex      sync.Mutex
	InProgress          bool
}

func (multiLobby *MultiplayerLobby) Join(client LobbyClient, password string) bool {
	if multiLobby.MatchInformation.GamePassword != password {
		return false
	}

	multiLobby.MatchInfoMutex.Lock()

	for i := 0; i < 7; i++ {
		if multiLobby.MatchInformation.SlotStatus[i] == packets.MultiplayerMatchSlotStatusOpen {
			for n := 0; n != 8; n++ {
				if multiLobby.MultiClients[n] != nil {
					packets.BanchoSendOsuUpdate(multiLobby.MultiClients[n].GetPacketQueue(), client.GetRelevantUserStats(), client.GetStatus())
				}
			}

			multiLobby.SetSlot(int32(i), client)
			multiLobby.MultiChannel.Join(client)

			multiLobby.MatchInfoMutex.Unlock()

			multiLobby.UpdateMatch()

			return true
		}
	}

	multiLobby.MatchInfoMutex.Unlock()

	return false
}

func (multiLobby *MultiplayerLobby) SetSlot(slot int32, client LobbyClient) {
	if client != nil {
		multiLobby.MatchInformation.SlotUserId[slot] = client.GetUserId()
		multiLobby.MatchInformation.SlotStatus[slot] = packets.MultiplayerMatchSlotStatusNotReady
		multiLobby.MultiClients[slot] = client

		if multiLobby.MatchInformation.MatchTeamType == packets.MultiplayerMatchTypeTagTeamVs || multiLobby.MatchInformation.MatchTeamType == packets.MultiplayerMatchTypeTeamVs {
			if slot%2 == 0 {
				multiLobby.MatchInformation.SlotTeam[slot] = packets.MultiplayerSlotTeamRed
			} else {
				multiLobby.MatchInformation.SlotTeam[slot] = packets.MultiplayerSlotTeamBlue
			}
		}
	} else {
		multiLobby.MatchInformation.SlotUserId[slot] = -1

		if multiLobby.MatchInformation.SlotStatus[slot] != packets.MultiplayerMatchSlotStatusLocked {
			multiLobby.MatchInformation.SlotStatus[slot] = packets.MultiplayerMatchSlotStatusOpen
		}

		multiLobby.MatchInformation.SlotTeam[slot] = packets.MultiplayerSlotTeamNeutral
		multiLobby.MultiClients[slot] = nil
	}
}

func (multiLobby *MultiplayerLobby) MoveSlot(oldSlot int, newSlot int) {
	if oldSlot == newSlot {
		return
	}

	currentStatus := multiLobby.MatchInformation.SlotStatus[oldSlot]

	multiLobby.SetSlot(int32(newSlot), multiLobby.MultiClients[oldSlot])
	multiLobby.SetSlot(int32(oldSlot), nil)

	multiLobby.MatchInformation.SlotStatus[newSlot] = currentStatus
}

func (multiLobby *MultiplayerLobby) UpdateMatch() {
	for i := 0; i != 8; i++ {
		if multiLobby.MultiClients[i] != nil {
			packets.BanchoSendMatchUpdate(multiLobby.MultiClients[i].GetPacketQueue(), multiLobby.MatchInformation)
		}
	}

	//Distribute in multiLobby as well
	BroadcastToLobby(func(packetQueue chan packets.BanchoPacket) {
		packets.BanchoSendMatchUpdate(packetQueue, multiLobby.MatchInformation)
	})
}

func (multiLobby *MultiplayerLobby) Part(client LobbyClient) {
	multiLobby.MatchInfoMutex.Lock()

	slot := multiLobby.GetSlotFromUserId(client.GetUserId())

	if slot == -1 {
		return
	}

	multiLobby.SetSlot(int32(slot), nil)

	if multiLobby.MatchHost == client {
		multiLobby.HandleHostLeave(slot)
	}

	packets.BanchoSendChannelRevoked(client.GetPacketQueue(), "#multiplayer")
	multiLobby.MultiChannel.Leave(client)

	if multiLobby.GetUsedUpSlots() == 0 {
		multiLobby.Disband()
	}

	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()
}

func (multiLobby *MultiplayerLobby) Disband() {
	RemoveMultiMatch(multiLobby.MatchInformation.MatchId)
}

func (multiLobby *MultiplayerLobby) HandleHostLeave(slot int) {
	if multiLobby.GetUsedUpSlots() == 0 {
		multiLobby.Disband()
	}

	for i := 0; i != 8; i++ {
		if multiLobby.MultiClients[i] != nil {
			multiLobby.MatchHost = multiLobby.MultiClients[i]

			if multiLobby.InProgress == false {
				multiLobby.MoveSlot(i, slot)
			}

			packets.BanchoSendMatchTransferHost(multiLobby.MatchHost.GetPacketQueue())
			multiLobby.MatchInformation.HostId = multiLobby.MatchHost.GetUserId()
		}
	}

	multiLobby.UpdateMatch()
}

func (multiLobby *MultiplayerLobby) TryChangeSlot(client LobbyClient, slotId int) {
	multiLobby.MatchInfoMutex.Lock()

	if multiLobby.MatchInformation.SlotStatus[slotId] == packets.MultiplayerMatchSlotStatusLocked || (multiLobby.MatchInformation.SlotStatus[slotId]&packets.MultiplayerMatchSlotStatusHasPlayer) > 0 {
		return
	}

	multiLobby.MoveSlot(multiLobby.GetSlotFromUserId(client.GetUserId()), slotId)
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()
}

func (multiLobby *MultiplayerLobby) ChangeTeam(client LobbyClient) {
	multiLobby.MatchInfoMutex.Lock()

	clientSlot := multiLobby.GetSlotFromUserId(client.GetUserId())

	if clientSlot == -1 {
		return
	}

	if multiLobby.MatchInformation.SlotTeam[clientSlot] == packets.MultiplayerSlotTeamRed {
		multiLobby.MatchInformation.SlotTeam[clientSlot] = packets.MultiplayerSlotTeamBlue
	} else if multiLobby.MatchInformation.SlotTeam[clientSlot] == packets.MultiplayerSlotTeamBlue {
		multiLobby.MatchInformation.SlotTeam[clientSlot] = packets.MultiplayerSlotTeamRed
	} else {
		multiLobby.MatchInformation.SlotTeam[clientSlot] = packets.MultiplayerSlotTeamRed
	}

	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()
}

func (multiLobby *MultiplayerLobby) TransferHost(client LobbyClient, slotId int) {
	multiLobby.MatchInfoMutex.Lock()

	if multiLobby.MatchHost != client {
		return
	}

	multiLobby.MatchHost = multiLobby.MultiClients[slotId]
	multiLobby.MatchInformation.HostId = multiLobby.MatchHost.GetUserId()

	packets.BanchoSendMatchTransferHost(multiLobby.MatchHost.GetPacketQueue())

	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()
}

func (multiLobby *MultiplayerLobby) ReadyUp(client LobbyClient) {
	multiLobby.MatchInfoMutex.Lock()

	clientSlot := multiLobby.GetSlotFromUserId(client.GetUserId())

	if clientSlot == -1 {
		return
	}

	multiLobby.MatchInformation.SlotStatus[clientSlot] = packets.MultiplayerMatchSlotStatusReady
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()
}

func (multiLobby *MultiplayerLobby) Unready(client LobbyClient) {
	multiLobby.MatchInfoMutex.Lock()

	clientSlot := multiLobby.GetSlotFromUserId(client.GetUserId())

	if clientSlot == -1 {
		return
	}

	multiLobby.MatchInformation.SlotStatus[clientSlot] = packets.MultiplayerMatchSlotStatusNotReady
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()
}

func (multiLobby *MultiplayerLobby) ChangeSettings(client LobbyClient, matchSettings packets.MultiplayerMatch) {
	multiLobby.MatchInfoMutex.Lock()

	if multiLobby.MatchHost != client {
		return
	}

	multiLobby.MatchInformation = matchSettings
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()
}

func (multiLobby *MultiplayerLobby) ChangeMods(client LobbyClient, newMods int32) {
	multiLobby.MatchInfoMutex.Lock()

	if multiLobby.MatchHost != client {
		return
	}

	multiLobby.MatchInformation.ActiveMods = uint16(newMods)
	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()
}

func (multiLobby *MultiplayerLobby) LockSlot(client LobbyClient, slotId int) {
	multiLobby.MatchInfoMutex.Lock()

	if multiLobby.MatchHost != client {
		return
	}

	if multiLobby.MultiClients[slotId] == multiLobby.MatchHost {
		return
	}

	if (multiLobby.MatchInformation.SlotStatus[slotId] & packets.MultiplayerMatchSlotStatusHasPlayer) > 0 {
		droppedClient := multiLobby.MultiClients[slotId]

		multiLobby.MatchInfoMutex.Unlock()

		droppedClient.LeaveCurrentMatch()
		packets.BanchoSendMatchUpdate(droppedClient.GetPacketQueue(), multiLobby.MatchInformation)

		multiLobby.MatchInfoMutex.Lock()
	}

	if multiLobby.MatchInformation.SlotStatus[slotId] == packets.MultiplayerMatchSlotStatusLocked {
		multiLobby.MatchInformation.SlotStatus[slotId] = packets.MultiplayerMatchSlotStatusOpen

		multiLobby.UpdateMatch()
		multiLobby.MatchInfoMutex.Unlock()

		return
	}

	if multiLobby.GetOpenSlotCount() > 2 && multiLobby.MatchInformation.SlotStatus[slotId] == packets.MultiplayerMatchSlotStatusOpen {
		multiLobby.MatchInformation.SlotStatus[slotId] = packets.MultiplayerMatchSlotStatusLocked
	}

	multiLobby.UpdateMatch()

	multiLobby.MatchInfoMutex.Unlock()
}
