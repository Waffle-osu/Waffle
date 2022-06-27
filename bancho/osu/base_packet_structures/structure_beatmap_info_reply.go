package base_packet_structures

import (
	"encoding/binary"
	"io"
)

type BeatmapInfoReply struct {
	BeatmapInfos []BeatmapInfo
}

func ReadBeatmapInfoReply(reader io.Reader) BeatmapInfoReply {
	infoReply := BeatmapInfoReply{}

	count := int32(0)

	binary.Read(reader, binary.LittleEndian, &count)

	for i := 0; i != int(count); i++ {
		infoReply.BeatmapInfos = append(infoReply.BeatmapInfos, ReadBeatmapInfo(reader))
	}

	return infoReply
}

func (infoReply BeatmapInfoReply) Write(writer io.Writer) {
	binary.Write(writer, binary.LittleEndian, int32(len(infoReply.BeatmapInfos)))

	for i := 0; i != len(infoReply.BeatmapInfos); i++ {
		infoReply.BeatmapInfos[i].Write(writer)
	}
}
