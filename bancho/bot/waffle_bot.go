package bot

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/lobby"
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/database"
	"sync"
	"time"
)

type WaffleBot struct {
	continueRunning bool

	logonTime time.Time

	lastReceive time.Time
	lastPing    time.Time

	joinedChannels map[string]*chat.Channel
	awayMessage    string

	spectators       map[int32]client_manager.WaffleClient
	spectatorMutex   sync.Mutex
	spectatingClient client_manager.WaffleClient

	isInLobby         bool
	currentMultiLobby *lobby.MultiplayerLobby

	UserData database.User
	Status   base_packet_structures.StatusUpdate

	OsuStats   database.UserStats
	TaikoStats database.UserStats
	CatchStats database.UserStats
	ManiaStats database.UserStats
}
