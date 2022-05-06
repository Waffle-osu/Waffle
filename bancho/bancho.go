package bancho

import (
	"Waffle/bancho/clients"
	"Waffle/helpers"
	"fmt"
	"net"
)

func RunBancho() {
	fmt.Printf("Running Bancho on 127.0.0.1:13381\n")

	//Creates the TCP server under which Waffle runs
	listener, err := net.Listen("tcp", "127.0.0.1:13381")

	if err != nil {
		helpers.Logger.Fatalf("[Bancho] Failed to Create TCP Server on 127.0.0.1:13381")
	}

	for {
		//Accept connections
		conn, err := listener.Accept()
		helpers.Logger.Printf("[Bancho] Connection Accepted!\n")

		if err != nil {
			continue
		}

		//Handle new connection
		go clients.HandleNewClient(conn)
	}
}
