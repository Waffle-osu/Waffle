package clients

import (
	"Waffle/waffle/chat"
	"Waffle/waffle/client_manager"
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

	joinedChannels []*chat.Channel

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
	client_manager.UnregisterClient(client)
	client_manager.BroadcastPacket(func(packetQueue chan packets.BanchoPacket) {
		packets.BanchoSendHandleOsuQuit(packetQueue, int32(client.UserData.UserID))
	})

	for _, channel := range client.joinedChannels {
		channel.Leave(client)
	}

	client.continueRunning = false

	client.connection.Close()
}
