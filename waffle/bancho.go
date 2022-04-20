package waffle

import (
	"fmt"
	"net"
)

type Bancho struct {
	server  net.Listener
	clients []Client
}

func CreateBancho() *Bancho {
	bancho := new(Bancho)

	listener, err := net.Listen("tcp", "127.0.0.1:13381")

	if err != nil {
		fmt.Printf("Failed to Create TCP Server on 127.0.0.1:13381")
	}

	bancho.server = listener
	bancho.clients = make([]Client, 128)

	return bancho
}

func (bancho Bancho) RunBancho() {
	fmt.Printf("Running Bancho on 127.0.0.1:13381\n")

	for {
		conn, err := bancho.server.Accept()
		fmt.Printf("Connection Accepted!\n")

		if err != nil {
			continue
		}

		go HandleNewClient(&bancho, conn)
	}
}
