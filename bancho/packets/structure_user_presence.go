package packets

import (
	"encoding/binary"
	"io"
)

type UserPresence struct {
	UserId          int32
	Username        string
	AvatarExtension uint8
	Timezone        uint8
	Country         uint8
	City            string
	Permissions     uint8
	Longitude       float32
	Latitude        float32
	Rank            int32
}

func ReadUserPresence(reader io.Reader) UserPresence {
	presence := UserPresence{}

	binary.Read(reader, binary.LittleEndian, &presence.UserId)
	presence.Username = string(ReadBanchoString(reader))
	binary.Read(reader, binary.LittleEndian, &presence.AvatarExtension)
	binary.Read(reader, binary.LittleEndian, &presence.Timezone)
	binary.Read(reader, binary.LittleEndian, &presence.Country)
	presence.City = string(ReadBanchoString(reader))
	binary.Read(reader, binary.LittleEndian, &presence.Permissions)
	binary.Read(reader, binary.LittleEndian, &presence.Longitude)
	binary.Read(reader, binary.LittleEndian, &presence.Latitude)
	binary.Read(reader, binary.LittleEndian, &presence.Rank)

	return presence
}

func (presence *UserPresence) WriteUserPresence(writer io.Writer) {
	binary.Write(writer, binary.LittleEndian, presence.UserId)
	binary.Write(writer, binary.LittleEndian, WriteBanchoString(presence.Username))
	binary.Write(writer, binary.LittleEndian, presence.AvatarExtension)
	binary.Write(writer, binary.LittleEndian, presence.Timezone)
	binary.Write(writer, binary.LittleEndian, presence.Country)
	binary.Write(writer, binary.LittleEndian, WriteBanchoString(presence.City))
	binary.Write(writer, binary.LittleEndian, presence.Permissions)
	binary.Write(writer, binary.LittleEndian, presence.Longitude)
	binary.Write(writer, binary.LittleEndian, presence.Latitude)
	binary.Write(writer, binary.LittleEndian, presence.Rank)
}
