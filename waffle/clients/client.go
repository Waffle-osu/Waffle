package clients

import (
	"Waffle/waffle/chat"
	"Waffle/waffle/database"
	"Waffle/waffle/packets"
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
	continueRunning bool

	lastReceive time.Time
	lastPing    time.Time

	joinedChannels []string

	PacketQueue chan packets.BanchoPacket

	UserData   database.User
	ClientData ClientInformation
	Status     packets.OsuStatus
	OsuStats   database.UserStats
	TaikoStats database.UserStats
	CatchStats database.UserStats
	ManiaStats database.UserStats
}

func CleanupClient(client *Client) {
	UnregisterClient(client)

	for _, channel := range client.joinedChannels {
		chat.LeaveChannel(client, channel)
	}

	client.continueRunning = false
	close(client.PacketQueue)

	client.connection.Close()
}
