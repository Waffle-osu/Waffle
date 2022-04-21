package packets

import (
	"Waffle/waffle/database"
	"bytes"
	"encoding/binary"
)

const (
	PresenceAvatarExtensionNone int8 = 0
	PresenceAvatarExtensionPng  int8 = 1
	PresenceAvatarExtensionJpg  int8 = 2
)

func BanchoSendUserPresence(packetQueue chan BanchoPacket, user database.User, stats database.UserStats, timezone int32) {
	buf := new(bytes.Buffer)

	//Write Data
	binary.Write(buf, binary.LittleEndian, int32(user.UserID))
	binary.Write(buf, binary.LittleEndian, WriteBanchoString(user.Username))
	binary.Write(buf, binary.LittleEndian, PresenceAvatarExtensionPng)
	binary.Write(buf, binary.LittleEndian, int8(timezone))
	binary.Write(buf, binary.LittleEndian, int8(user.Country))
	binary.Write(buf, binary.LittleEndian, WriteBanchoString("No city for you!"))
	binary.Write(buf, binary.LittleEndian, int8(user.Privileges&0b11111111))
	binary.Write(buf, binary.LittleEndian, float32(0.0))
	binary.Write(buf, binary.LittleEndian, float32(0.0))
	binary.Write(buf, binary.LittleEndian, int32(1)) //TODO: rank

	packetBytes := buf.Bytes()
	packetLength := len(packetBytes)

	packet := BanchoPacket{
		PacketId:          BanchoUserPresence,
		PacketCompression: 0,
		PacketSize:        int32(packetLength),
		PacketData:        packetBytes,
	}

	packetQueue <- packet
}
