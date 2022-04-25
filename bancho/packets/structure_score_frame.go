package packets

import (
	"encoding/binary"
	"io"
)

type ScoreFrame struct {
	Time         int32
	Id           uint8
	Count300     uint16
	Count100     uint16
	Count50      uint16
	CountGeki    uint16
	CountKatu    uint16
	CountMiss    uint16
	TotalScore   int32
	MaxCombo     uint16
	CurrentCombo uint16
	Perfect      bool
	CurrentHp    uint8
	TagByte      uint8
}

func ReadScoreFrame(reader io.Reader) ScoreFrame {
	scoreFrame := ScoreFrame{}

	binary.Read(reader, binary.LittleEndian, &scoreFrame.Time)
	binary.Read(reader, binary.LittleEndian, &scoreFrame.Id)
	binary.Read(reader, binary.LittleEndian, &scoreFrame.Count300)
	binary.Read(reader, binary.LittleEndian, &scoreFrame.Count100)
	binary.Read(reader, binary.LittleEndian, &scoreFrame.Count50)
	binary.Read(reader, binary.LittleEndian, &scoreFrame.CountGeki)
	binary.Read(reader, binary.LittleEndian, &scoreFrame.CountKatu)
	binary.Read(reader, binary.LittleEndian, &scoreFrame.CountMiss)
	binary.Read(reader, binary.LittleEndian, &scoreFrame.TotalScore)
	binary.Read(reader, binary.LittleEndian, &scoreFrame.MaxCombo)
	binary.Read(reader, binary.LittleEndian, &scoreFrame.CurrentCombo)
	binary.Read(reader, binary.LittleEndian, &scoreFrame.Perfect)
	binary.Read(reader, binary.LittleEndian, &scoreFrame.CurrentHp)
	binary.Read(reader, binary.LittleEndian, &scoreFrame.TagByte)

	return scoreFrame
}

func (scoreFrame *ScoreFrame) WriteScoreFrame(writer io.Writer) {
	binary.Write(writer, binary.LittleEndian, scoreFrame.Time)
	binary.Write(writer, binary.LittleEndian, scoreFrame.Id)
	binary.Write(writer, binary.LittleEndian, scoreFrame.Count300)
	binary.Write(writer, binary.LittleEndian, scoreFrame.Count100)
	binary.Write(writer, binary.LittleEndian, scoreFrame.Count50)
	binary.Write(writer, binary.LittleEndian, scoreFrame.CountGeki)
	binary.Write(writer, binary.LittleEndian, scoreFrame.CountKatu)
	binary.Write(writer, binary.LittleEndian, scoreFrame.CountMiss)
	binary.Write(writer, binary.LittleEndian, scoreFrame.TotalScore)
	binary.Write(writer, binary.LittleEndian, scoreFrame.MaxCombo)
	binary.Write(writer, binary.LittleEndian, scoreFrame.CurrentCombo)
	binary.Write(writer, binary.LittleEndian, scoreFrame.Perfect)
	binary.Write(writer, binary.LittleEndian, scoreFrame.CurrentHp)
	binary.Write(writer, binary.LittleEndian, scoreFrame.TagByte)
}
