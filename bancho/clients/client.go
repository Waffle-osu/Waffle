package clients

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/lobby"
	"Waffle/bancho/packets"
	"Waffle/database"
	"Waffle/logger"
	"net"
	"sync"
	"time"
)

const (
	// ReceiveTimeout 48 Seconds
	ReceiveTimeout = 16
	// PingTimeout 8 Seconds
	PingTimeout = 8
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

	clean      bool
	cleanMutex sync.Mutex

	joinedChannels map[string]*chat.Channel
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

// CleanupClient cleans the client up, leaves spectator and the lobby and the multi match if applicable, also lets everyone know its departure
func (client *Client) CleanupClient(reason string) {
	client.cleanMutex.Lock()

	if client.clean {
		return
	}

	logger.Logger.Printf("[Bancho@Client] Cleaning up %s; Reason: %s\n", client.UserData.Username, reason)

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

	client.connection.Close()

	client.clean = true

	client.cleanMutex.Unlock()
}

// Cut cuts the client's connection and forces a disconnect.
func (client *Client) Cut() {
	client.continueRunning = false
	client.connection.Close()
}
