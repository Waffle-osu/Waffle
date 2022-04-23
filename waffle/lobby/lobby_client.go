package lobby

import (
	"Waffle/waffle/database"
	"Waffle/waffle/packets"
)

type LobbyClient interface {
	GetPacketQueue() chan packets.BanchoPacket
	GetUserId() int32
	GetUserData() database.User
}
