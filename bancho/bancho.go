package bancho

import (
	"Waffle/bancho/clients"
	"Waffle/config"
	"Waffle/helpers"
	"net"
)

func RunBancho() {
	//Creates the TCP server under which Waffle runs
	listener, err := net.Listen("tcp", config.BanchoIp)

	if err != nil {
		helpers.Logger.Fatalf("[Bancho] Failed to Create TCP Server on %s", config.BanchoIp)
	}

	helpers.Logger.Printf("Running Bancho on %s\n", config.BanchoIp)

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
