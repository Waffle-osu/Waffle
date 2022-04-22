package client_manager

import (
	"Waffle/waffle/database"
	"Waffle/waffle/packets"
)

type OsuClient interface {
	GetUserId() int32
	GetPacketQueue() chan packets.BanchoPacket
	GetRelevantUserStats() database.UserStats
	GetUserStatus() packets.OsuStatus
	GetUserData() database.User
	GetClientTimezone() int32
}
