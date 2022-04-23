package lobby

import "Waffle/waffle/packets"

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
		if (multiLobby.MatchInformation.SlotStatus[i] & packets.MultiplayerMatchSlotStatusHasPlayer) > 0 {
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
