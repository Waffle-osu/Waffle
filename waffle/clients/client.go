package clients

import (
	"Waffle/waffle/database"
	"Waffle/waffle/packets"
	"bufio"
	"net"
	"time"
)

const (
	// ReceiveTimeout 15 Seconds
	ReceiveTimeout = 15000000000
	// PingTimeout 10 Seconds
	PingTimeout = 10000000000
)

type ClientInformation struct {
	Timezone       int32
	Version        string
	AllowCity      bool
	OsuClientHash  string
	MacAddressHash string
}

type Client struct {
	connection      net.Conn
	bufReader       *bufio.Reader
	continueRunning bool

	lastRecieve time.Time
	lastPing    time.Time

	PacketQueue chan packets.BanchoPacket

	UserData   database.User
	ClientData ClientInformation
	Status     packets.OsuStatus
	OsuStats   database.UserStats
	TaikoStats database.UserStats
	CatchStats database.UserStats
	ManiaStats database.UserStats
}
