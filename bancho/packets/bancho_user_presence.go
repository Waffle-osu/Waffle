package packets

import (
	"Waffle/database"
	"bytes"
)

const (
	PresenceAvatarExtensionNone uint8 = 0
	PresenceAvatarExtensionPng  uint8 = 1
	PresenceAvatarExtensionJpg  uint8 = 2
)

func BanchoSendUserPresence(packetQueue chan BanchoPacket, user database.User, stats database.UserStats, timezone int32) {
	buf := new(bytes.Buffer)

	presence := UserPresence{
		UserId:          int32(user.UserID),
		Username:        user.Username,
		AvatarExtension: PresenceAvatarExtensionPng,
		Timezone:        uint8(timezone),
		Country:         uint8(user.Country),
		City:            "",
		Permissions:     uint8(user.Privileges),
		Longitude:       0,
		Latitude:        0,
		Rank:            int32(stats.Rank),
	}

	presence.WriteUserPresence(buf)

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
