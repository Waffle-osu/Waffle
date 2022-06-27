package base_packet_structures

import (
	"Waffle/helpers"
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
	matchJoin.Password = string(helpers.ReadBanchoString(reader))

	return matchJoin
}

func (matchJoin *MatchJoin) WriteMatchJoin(writer io.Writer) {
	binary.Write(writer, binary.LittleEndian, matchJoin.MatchId)
	binary.Write(writer, binary.LittleEndian, helpers.WriteBanchoString(matchJoin.Password))
}
