package base_packet_structures

import (
	"Waffle/helpers/serialization"
	"encoding/binary"
	"io"
)

type BeatmapInfo struct {
	InfoId          int16
	BeatmapId       int32
	BeatmapSetId    int32
	ThreadId        int32
	Ranked          uint8
	OsuRank         uint8
	TaikoRank       uint8
	CatchRank       uint8
	BeatmapChecksum string
}

func ReadBeatmapInfo(reader io.Reader) BeatmapInfo {
	beatmapInfo := BeatmapInfo{}

	binary.Read(reader, binary.LittleEndian, &beatmapInfo.InfoId)
	binary.Read(reader, binary.LittleEndian, &beatmapInfo.BeatmapId)
	binary.Read(reader, binary.LittleEndian, &beatmapInfo.BeatmapSetId)
	binary.Read(reader, binary.LittleEndian, &beatmapInfo.ThreadId)
	binary.Read(reader, binary.LittleEndian, &beatmapInfo.Ranked)
	binary.Read(reader, binary.LittleEndian, &beatmapInfo.OsuRank)
	binary.Read(reader, binary.LittleEndian, &beatmapInfo.CatchRank)
	binary.Read(reader, binary.LittleEndian, &beatmapInfo.TaikoRank)
	beatmapInfo.BeatmapChecksum = string(serialization.ReadBanchoString(reader))

	return beatmapInfo
}

func (beatmapInfo BeatmapInfo) Write(writer io.Writer) {
	binary.Write(writer, binary.LittleEndian, beatmapInfo.InfoId)
	binary.Write(writer, binary.LittleEndian, beatmapInfo.BeatmapId)
	binary.Write(writer, binary.LittleEndian, beatmapInfo.BeatmapSetId)
	binary.Write(writer, binary.LittleEndian, beatmapInfo.ThreadId)
	binary.Write(writer, binary.LittleEndian, beatmapInfo.Ranked)
	binary.Write(writer, binary.LittleEndian, beatmapInfo.OsuRank)
	binary.Write(writer, binary.LittleEndian, beatmapInfo.CatchRank)
	binary.Write(writer, binary.LittleEndian, beatmapInfo.TaikoRank)
	binary.Write(writer, binary.LittleEndian, serialization.WriteBanchoString(beatmapInfo.BeatmapChecksum))
}
