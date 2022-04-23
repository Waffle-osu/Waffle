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

func CreateNewMatch(match packets.MultiplayerMatch, host LobbyClient) *MultiplayerLobby {
	multiLobby := new(MultiplayerLobby)

	multiLobby.MultiChannel = new(chat.Channel)
	multiLobby.MultiChannel.Name = "#multiplayer"
	multiLobby.MultiChannel.Description = ""
	multiLobby.MultiChannel.AdminWrite = false
	multiLobby.MultiChannel.AdminRead = false
	multiLobby.MultiChannel.Autojoin = false
	multiLobby.MultiChannel.Clients = []chat.ChatClient{}
	multiLobby.MultiChannel.ClientMutex = sync.Mutex{}

	multiLobby.MatchInformation = match
	multiLobby.MatchHost = host

	return multiLobby
}

func (multiLobby *MultiplayerLobby) Join(client LobbyClient) bool {
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

func (multiLobby *MultiplayerLobby) GetSlotFromUserId(userId int32) int {
	for i := 0; i != 8; i++ {
		if multiLobby.MatchInformation.SlotUserId[i] == userId {
			return i
		}
	}

	return -1
}

func (multiLobby *MultiplayerLobby) GetOpenSlotCount() int {
	count := 0

	for i := 0; i != 8; i++ {
		if multiLobby.MatchInformation.SlotStatus[i] == packets.MultiplayerMatchSlotStatusLocked {
			count++
		}
	}

	return count
}

func (multiLobby *MultiplayerLobby) HaveAllPlayersSkipped() bool {
	for i := 0; i != 8; i++ {
		if multiLobby.MatchInformation.SlotStatus[i] == packets.MultiplayerMatchSlotStatusPlaying && multiLobby.PlayerSkipRequested[i] == false {
			return false
		}
	}

	return true
}

func (multiLobby *MultiplayerLobby) HaveAllPlayersCompleted() bool {
	count := 0

	for i := 0; i != 8; i++ {
		if multiLobby.PlayerCompleted[i] == true {
			count++
		}
	}

	return count == multiLobby.GetUsedUpSlots()
}

func (multiLobby *MultiplayerLobby) HaveAllPlayersLoaded() bool {
	count := 0

	for i := 0; i != 8; i++ {
		if multiLobby.PlayersLoaded[i] == true {
			count++
		}
	}

	return count == multiLobby.GetUsedUpSlots()
}

func (multiLobby *MultiplayerLobby) GetUsedUpSlots() int {
	count := 0

	for i := 0; i != 8; i++ {
		if multiLobby.MatchInformation.SlotStatus[i] == packets.MultiplayerMatchSlotStatusHasPlayer {
			count++
		}
	}

	return count
}

func (multiLobby *MultiplayerLobby) HaveAllPlayersFinished() bool {
	finished := 0

	for i := 0; i != 8; i++ {
		if multiLobby.PlayerCompleted[i] == true {
			finished++
		}
	}

	return finished == multiLobby.GetUsedUpSlots()
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
	//if nobody is inside: disband and return

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
