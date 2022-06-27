package base_packet_structures

import (
	"Waffle/helpers/serialization"
	"encoding/binary"
	"io"
)

type MatchJoin struct {
	MatchId  int32
	Password string
}

func ReadMatchJoin(reader io.Reader) MatchJoin {
	matchJoin := MatchJoin{}

	binary.Read(reader, binary.LittleEndian, &matchJoin.MatchId)
	matchJoin.Password = string(serialization.ReadBanchoString(reader))

	return matchJoin
}

func (matchJoin MatchJoin) Write(writer io.Writer) {
	binary.Write(writer, binary.LittleEndian, matchJoin.MatchId)
	binary.Write(writer, binary.LittleEndian, serialization.WriteBanchoString(matchJoin.Password))
}
