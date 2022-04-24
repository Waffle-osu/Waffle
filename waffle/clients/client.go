package clients

import (
	"Waffle/waffle/chat"
	"Waffle/waffle/client_manager"
	"Waffle/waffle/database"
	"Waffle/waffle/lobby"
	"Waffle/waffle/packets"
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	// ReceiveTimeout 48 Seconds
	ReceiveTimeout = 48000000000
	// PingTimeout 8 Seconds
	PingTimeout = 8000000000
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
	awayMessage    string

	spectators       map[int32]client_manager.OsuClient
	spectatorMutex   sync.Mutex
	spectatingClient client_manager.OsuClient

	isInLobby         bool
	currentMultiLobby *lobby.MultiplayerLobby

	PacketQueue chan packets.BanchoPacket

	UserData    database.User
	ClientData  ClientInformation
	Status      packets.StatusUpdate
	OsuStats    database.UserStats
	TaikoStats  database.UserStats
	CatchStats  database.UserStats
	ManiaStats  database.UserStats
	FriendsList []database.FriendEntry
}

func (client *Client) CleanupClient() {
	fmt.Printf("Cleaning up %s\n", client.UserData.Username)

	if client.spectatingClient != nil {
		client.spectatingClient.InformSpectatorLeft(client)
	}

	client.spectators = map[int32]client_manager.OsuClient{}

	if client.isInLobby {
		lobby.PartLobby(client)
	}

	if client.currentMultiLobby != nil {
		client.LeaveCurrentMatch()
	}

	client_manager.UnregisterClient(client)
	client_manager.BroadcastPacket(func(packetQueue chan packets.BanchoPacket) {
		packets.BanchoSendHandleOsuQuit(packetQueue, int32(client.UserData.UserID))
	})

	for _, channel := range client.joinedChannels {
		channel.Leave(client)
	}
}

func (client *Client) Cut() {
	client.continueRunning = false
	client.connection.Close()
}
