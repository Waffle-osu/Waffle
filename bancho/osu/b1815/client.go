package b1815

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/lobby"
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/database"
	"Waffle/helpers"
	"context"
	"net"
	"sync"
	"time"
)

const (
	// ReceiveTimeout 30 Seconds
	ReceiveTimeout = 30
	// PingTimeout 6 Seconds, after 5 Pings we disconnect the client.
	PingTimeout = 6
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

	logonTime time.Time

	lastReceive time.Time
	lastPing    time.Time

	clean      bool
	cleanMutex sync.Mutex

	joinedChannels map[string]*chat.Channel
	awayMessage    string

	spectators       map[int32]client_manager.WaffleClient
	spectatorMutex   sync.Mutex
	spectatingClient client_manager.WaffleClient

	isInLobby         bool
	currentMultiLobby *lobby.MultiplayerLobby

	packetListeners []PacketEvent

	waffleGuardContext WaffleGuardContext

	maintainCancel context.CancelFunc

	PacketQueue chan []byte

	UserData    database.User
	ClientData  ClientInformation
	Status      base_packet_structures.StatusUpdate
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

	helpers.Logger.Printf("[Bancho@Client] Cleaning up %s; Reason: %s\n", client.UserData.Username, reason)

	if client.spectatingClient != nil {
		client.spectatingClient.BanchoSpectatorLeft(client.GetUserId())
	}

	client.spectators = map[int32]client_manager.WaffleClient{}

	if client.isInLobby {
		lobby.PartLobby(client)
	}

	if client.currentMultiLobby != nil {
		client.LeaveCurrentMatch()
	}

	client_manager.UnregisterClient(client)
	client_manager.BroadcastPacketOsu(func(broadcastClient client_manager.WaffleClient) {
		broadcastClient.BanchoHandleOsuQuit(client.GetUserId(), client.GetUsername())
	})

	for _, channel := range client.joinedChannels {
		channel.Leave(client)
	}

	client.connection.Close()
	client.continueRunning = false

	client.clean = true
	client.cleanMutex.Unlock()

	//Cancels the outgoing packet queue, and pingers
	client.maintainCancel()
}

// Cut cuts the client's connection and forces a disconnect.
func (client *Client) Cut() {
	client.continueRunning = false
	client.connection.Close()
}
