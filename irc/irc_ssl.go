package irc

import (
	"Waffle/helpers"
	"net"
)

func RunIrcSSL() {
	helpers.Logger.Printf("Running IRC/SSL on 127.0.0.1:6697\n")

	listener, err := net.Listen("tcp", "127.0.0.1:6697")

	if err != nil {
		helpers.Logger.Printf("[IRC/SSL] Failed to create TCP Listener for IRC/SSL on 127.0.0.1:6697\n")
	}

	for {
		conn, err := listener.Accept()

		helpers.Logger.Printf("[IRC/SSL] Accepted Connection!\n")

		if err != nil {
			continue
		}

		
	}
}
