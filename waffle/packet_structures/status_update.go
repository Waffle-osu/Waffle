package packet_structures

import (
	"Waffle/waffle/packets"
	"encoding/binary"
	"io"
)

type StatusUpdate struct {
	Status          uint8
	StatusText      string
	BeatmapChecksum string
	CurrentMods     uint16
	Playmode        uint8
	BeatmapId       int32
}

func ReadStatusUpdate(reader io.Reader) StatusUpdate {
	statusUpdate := StatusUpdate{}

	binary.Read(reader, binary.LittleEndian, &statusUpdate.Status)
	statusUpdate.StatusText = string(packets.ReadBanchoString(reader))
	statusUpdate.BeatmapChecksum = string(packets.ReadBanchoString(reader))
	binary.Read(reader, binary.LittleEndian, &statusUpdate.CurrentMods)
	binary.Read(reader, binary.LittleEndian, &statusUpdate.Playmode)
	binary.Read(reader, binary.LittleEndian, &statusUpdate.BeatmapId)

	return statusUpdate
}

