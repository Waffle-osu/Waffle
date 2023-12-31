package base_packet_structures

type StatusUpdate struct {
	Status          uint8
	StatusText      string
	BeatmapChecksum string
	CurrentMods     uint16
	Playmode        uint8
	BeatmapId       int32
}
