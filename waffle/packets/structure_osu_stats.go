package packets

import (
	"encoding/binary"
	"io"
)

type OsuStats struct {
	UserId      int32
	Status      StatusUpdate
	RankedScore int64
	Accuracy    float32
	Playcount   int32
	TotalScore  int64
	Rank        int32
}

func ReadOsuStats(reader io.Reader) OsuStats {
	stats := OsuStats{}

	binary.Read(reader, binary.LittleEndian, &stats.UserId)
	stats.Status = ReadStatusUpdate(reader)
	binary.Read(reader, binary.LittleEndian, &stats.RankedScore)
	binary.Read(reader, binary.LittleEndian, &stats.Accuracy)
	binary.Read(reader, binary.LittleEndian, &stats.Playcount)
	binary.Read(reader, binary.LittleEndian, &stats.TotalScore)
	binary.Read(reader, binary.LittleEndian, &stats.Rank)

	return stats
}

func (stats *OsuStats) WriteOsuStats(writer io.Writer) {
	binary.Write(writer, binary.LittleEndian, stats.UserId)
	stats.Status.WriteStatusUpdate(writer)
	binary.Write(writer, binary.LittleEndian, stats.RankedScore)
	binary.Write(writer, binary.LittleEndian, stats.Accuracy)
	binary.Write(writer, binary.LittleEndian, stats.Playcount)
	binary.Write(writer, binary.LittleEndian, stats.TotalScore)
	binary.Write(writer, binary.LittleEndian, stats.Rank)
}
