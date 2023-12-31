package base_packet_structures

const (
	MultiplayerMatchSlotStatusOpen       uint8 = 1
	MultiplayerMatchSlotStatusLocked     uint8 = 2
	MultiplayerMatchSlotStatusNotReady   uint8 = 4
	MultiplayerMatchSlotStatusReady      uint8 = 8
	MultiplayerMatchSlotStatusMissingMap uint8 = 16
	MultiplayerMatchSlotStatusPlaying    uint8 = 32
	MultiplayerMatchSlotStatusCompleted  uint8 = 64
	MultiplayerMatchSlotStatusHasPlayer  uint8 = MultiplayerMatchSlotStatusNotReady | MultiplayerMatchSlotStatusReady | MultiplayerMatchSlotStatusMissingMap | MultiplayerMatchSlotStatusPlaying | MultiplayerMatchSlotStatusCompleted
	MultiplayerMatchSlotStatusQuit       uint8 = 128

	MultiplayerSlotTeamNeutral uint8 = 0
	MultiplayerSlotTeamBlue    uint8 = 1
	MultiplayerSlotTeamRed     uint8 = 2

	MultiplayerMatchTypeHeadToHead uint8 = 0
	MultiplayerMatchTypeTagCoop    uint8 = 1
	MultiplayerMatchTypeTeamVs     uint8 = 2
	MultiplayerMatchTypeTagTeamVs  uint8 = 3

	MultiplayerMatchScoreTypeScore    uint8 = 0
	MultiplayerMatchScoreTypeAccuracy uint8 = 1
)

type MultiplayerMatch struct {
	MatchId          uint16
	InProgress       bool
	MatchType        uint8
	ActiveMods       uint16
	GameName         string
	GamePassword     string
	BeatmapName      string
	BeatmapId        int32
	BeatmapChecksum  string
	SlotStatus       [8]uint8
	SlotTeam         [8]uint8
	SlotUserId       [8]int32 `multi:"SlotStatus" multiAnd:"124"` // 124 == MultiplayerMatchSlotStatusHasPlayer, for the special tag implementation: reflection_serializer.go:350 and :123
	HostId           int32
	Playmode         uint8
	MatchScoringType uint8
	MatchTeamType    uint8
}
