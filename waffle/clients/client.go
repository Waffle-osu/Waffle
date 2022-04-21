package clients

import (
	"Waffle/waffle/database"
	"Waffle/waffle/packets"
	"bufio"
	"container/list"
	"net"
	"sync"
)

type ClientInformation struct {
	Timezone       int32
	Version        string
	AllowCity      bool
	OsuClientHash  string
	MacAddressHash string
}

type Client struct {
	Connection       net.Conn
	ClientData       ClientInformation
	BufReader        *bufio.Reader
	PacketQueue      *list.List
	PacketQueueMutex sync.Mutex
	UserData         database.User
	Status           packets.OsuStatus
	OsuStats         database.UserStats
	TaikoStats       database.UserStats
	CatchStats       database.UserStats
	ManiaStats       database.UserStats
}
