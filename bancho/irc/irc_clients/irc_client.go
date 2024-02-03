package irc_clients

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/irc/irc_messages"
	"Waffle/bancho/lobby"
	"Waffle/database"
	"Waffle/helpers"
	"bufio"
	"context"
	"net"
	"sync"
	"time"
)

var MOTD []string = []string{"",
	" _       __      __________     __",
	"| |     / /___ _/ __/ __/ /__  / /",
	"| | /| / / __ `/ /_/ /_/ / _ \\/ / ",
	"| |/ |/ / /_/ / __/ __/ /  __/_/  ",
	"|__/|__/\\__,_/_/ /_/ /_/\\___(_)   ",
	"                                 ",
}

const (
	ReceiveTimeout = 30
	PingTimeout    = 6
)

type IrcClient struct {
	connection      net.Conn
	reader          *bufio.Reader
	continueRunning bool
	packetQueue     chan irc_messages.Message

	logonTime time.Time

	lastReceive time.Time
	lastPing    time.Time

	joinedChannels map[string]*chat.Channel
	awayMessage    string

	currentMultiLobby *lobby.MultiplayerLobby

	lastPingToken string

	clean      bool
	cleanMutex sync.Mutex

	maintainCancel context.CancelFunc

	//Name used to address you on IRC
	//Must be unique across the network
	//This is the username used in /kick commands and similar
	Nickname string

	//This is used to populate the real name field when using /whois
	//can contain most characters
	Realname string

	//Is mainly used for people using 1 computer for more than 1 IRC User
	//To differenciate between them.
	//Also cannot be changed without a reconnect
	Username string

	//Password provided by the IRC Client
	Password string

	UserData database.User

	//Is it a IRC osu! client
	IsOsu bool
	//Next message sent is the username
	IsAwaitingUsername bool
	//Next message sent is the OTP
	IsAwaitingOtp bool

	ClientVersion client_manager.ClientVersion
}

func (client *IrcClient) CleanupClient(reason string) {
	client.cleanMutex.Lock()

	if client.clean {
		return
	}

	helpers.Logger.Printf("[IRC@IrcClient] Cleaning up %s; Reason: %s", client.Username, reason)

	if client.currentMultiLobby != nil {
		client.currentMultiLobby.Part(client)
		client.currentMultiLobby = nil
	}

	client_manager.ClientManager.UnregisterClient(client)

	for _, channel := range client.joinedChannels {
		channel.Leave(client)
	}

	client.connection.Close()

	client.clean = true
	client.cleanMutex.Unlock()

	if client.maintainCancel != nil {
		client.maintainCancel()
	}

	client.continueRunning = false
}

func (client *IrcClient) Cut() {
	client.continueRunning = false
	client.connection.Close()
}
