package base_packet_structures

import (
	"Waffle/helpers/serialization"
	"encoding/binary"
	"io"
)

type BeatmapInfoRequest struct {
	Filenames  []string
	BeatmapIds []int32
}

func ReadBeatmapInfoRequest(reader io.Reader) BeatmapInfoRequest {
	infoRequest := BeatmapInfoRequest{}

	filenameCount := int32(0)

	binary.Read(reader, binary.LittleEndian, &filenameCount)

	for i := 0; i != int(filenameCount); i++ {
		infoRequest.Filenames = append(infoRequest.Filenames, string(serialization.ReadBanchoString(reader)))
	}

	idCount := int32(0)

	binary.Read(reader, binary.LittleEndian, &idCount)

	for i := 0; i != int(idCount); i++ {
		beatmapId := int32(0)

		binary.Read(reader, binary.LittleEndian, &beatmapId)

		infoRequest.BeatmapIds = append(infoRequest.BeatmapIds, beatmapId)
	}

	return infoRequest
}

func (infoRequest BeatmapInfoRequest) Write(writer io.Writer) {
	binary.Write(writer, binary.LittleEndian, int32(len(infoRequest.Filenames)))

	for i := 0; i != len(infoRequest.Filenames); i++ {
		binary.Write(writer, binary.LittleEndian, serialization.WriteBanchoString(infoRequest.Filenames[i]))
	}

	binary.Write(writer, binary.LittleEndian, int32(len(infoRequest.BeatmapIds)))

	for i := 0; i != len(infoRequest.BeatmapIds); i++ {
		binary.Write(writer, binary.LittleEndian, infoRequest.BeatmapIds[i])
	}
}
