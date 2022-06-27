package base_packet_structures

import (
	"Waffle/helpers"
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
	statusUpdate.StatusText = string(helpers.ReadBanchoString(reader))
	statusUpdate.BeatmapChecksum = string(helpers.ReadBanchoString(reader))
	binary.Read(reader, binary.LittleEndian, &statusUpdate.CurrentMods)
	binary.Read(reader, binary.LittleEndian, &statusUpdate.Playmode)
	binary.Read(reader, binary.LittleEndian, &statusUpdate.BeatmapId)

	return statusUpdate
}

func (statusUpdate *StatusUpdate) WriteStatusUpdate(writer io.Writer) {
	binary.Write(writer, binary.LittleEndian, statusUpdate.Status)
	binary.Write(writer, binary.LittleEndian, helpers.WriteBanchoString(statusUpdate.StatusText))
	binary.Write(writer, binary.LittleEndian, helpers.WriteBanchoString(statusUpdate.BeatmapChecksum))
	binary.Write(writer, binary.LittleEndian, statusUpdate.CurrentMods)
	binary.Write(writer, binary.LittleEndian, statusUpdate.Playmode)
	binary.Write(writer, binary.LittleEndian, statusUpdate.BeatmapId)
}
