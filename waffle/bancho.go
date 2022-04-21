package waffle

import (
	"Waffle/waffle/chat"
	"Waffle/waffle/clients"
	"Waffle/waffle/database"
	"fmt"
	"net"
)

type Bancho struct {
	Server         net.Listener
	WorkerChannels []chan struct{}
}

func CreateBancho() *Bancho {
	bancho := new(Bancho)

	chat.InitializeChannels()
	database.Initialize()

	listener, err := net.Listen("tcp", "127.0.0.1:13381")

	if err != nil {
		fmt.Printf("Failed to Create TCP Server on 127.0.0.1:13381")
	}

	bancho.Server = listener

	return bancho
}

func (bancho *Bancho) RunBancho() {
	fmt.Printf("Running Bancho on 127.0.0.1:13381\n")

	for {
		conn, err := bancho.Server.Accept()
		fmt.Printf("Connection Accepted!\n")

		if err != nil {
			continue
		}

		go clients.HandleNewClient(conn)
	}
}
