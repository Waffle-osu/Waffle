package packets

import (
	"encoding/binary"
	"io"
)

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
	ActiveMods       int16
	GameName         string
	GamePassword     string
	BeatmapName      string
	BeatmapId        int32
	BeatmapChecksum  string
	SlotStatus       [8]uint8
	SlotTeam         [8]uint8
	SlotUserId       [8]int32
	HostId           int32
	Playmode         uint8
	MatchScoringType uint8
	MatchTeamType    uint8
}

func ReadMultiplayerMatch(reader io.Reader) MultiplayerMatch {
	match := MultiplayerMatch{}

	binary.Read(reader, binary.LittleEndian, &match.MatchId)
	binary.Read(reader, binary.LittleEndian, &match.InProgress)
	binary.Read(reader, binary.LittleEndian, &match.MatchType)
	binary.Read(reader, binary.LittleEndian, &match.ActiveMods)
	match.GameName = string(ReadBanchoString(reader))
	match.GamePassword = string(ReadBanchoString(reader))
	match.BeatmapName = string(ReadBanchoString(reader))
	binary.Read(reader, binary.LittleEndian, &match.BeatmapId)
	match.BeatmapChecksum = string(ReadBanchoString(reader))
	binary.Read(reader, binary.LittleEndian, &match.SlotStatus)
	binary.Read(reader, binary.LittleEndian, &match.SlotTeam)

	for i := 0; i != 8; i++ {
		if (match.SlotStatus[i] & MultiplayerMatchSlotStatusHasPlayer) > 0 {
			binary.Read(reader, binary.LittleEndian, &match.SlotUserId[i])
		}
	}

	binary.Read(reader, binary.LittleEndian, &match.HostId)
	binary.Read(reader, binary.LittleEndian, &match.Playmode)
	binary.Read(reader, binary.LittleEndian, &match.MatchScoringType)
	binary.Read(reader, binary.LittleEndian, &match.MatchTeamType)

	return match
}

func (match *MultiplayerMatch) WriteMultiplayerMatch(writer io.Writer) {
	binary.Write(writer, binary.LittleEndian, match.MatchId)
	binary.Write(writer, binary.LittleEndian, match.InProgress)
	binary.Write(writer, binary.LittleEndian, match.MatchType)
	binary.Write(writer, binary.LittleEndian, match.ActiveMods)
	binary.Write(writer, binary.LittleEndian, WriteBanchoString(match.GameName))
	binary.Write(writer, binary.LittleEndian, WriteBanchoString(match.GamePassword))
	binary.Write(writer, binary.LittleEndian, WriteBanchoString(match.BeatmapName))
	binary.Write(writer, binary.LittleEndian, match.BeatmapId)
	binary.Write(writer, binary.LittleEndian, WriteBanchoString(match.BeatmapChecksum))
	binary.Write(writer, binary.LittleEndian, match.SlotStatus)
	binary.Write(writer, binary.LittleEndian, match.SlotTeam)

	for i := 0; i != 8; i++ {
		if (match.SlotStatus[i] & MultiplayerMatchSlotStatusHasPlayer) > 0 {
			binary.Write(writer, binary.LittleEndian, match.SlotUserId[i])
		}
	}

	binary.Write(writer, binary.LittleEndian, match.HostId)
	binary.Write(writer, binary.LittleEndian, match.Playmode)
	binary.Write(writer, binary.LittleEndian, match.MatchScoringType)
	binary.Write(writer, binary.LittleEndian, match.MatchTeamType)
}
