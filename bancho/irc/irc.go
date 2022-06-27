package irc

import (
	"Waffle/bancho/irc/irc_clients"
	"Waffle/config"
	"Waffle/helpers"
	"net"
)

func RunIrc() {
	if config.HostIrc == "false" {
		return
	}

	listener, err := net.Listen("tcp", config.IrcIp)

	if err != nil {
		helpers.Logger.Printf("[IRC] Failed to create TCP Listener for IRC on %s\n", config.IrcIp)
	}

	helpers.Logger.Printf("Running IRC on %s\n", config.IrcIp)

	for {
		conn, err := listener.Accept()

		helpers.Logger.Printf("[IRC] Accepted Connection!\n")

		if err != nil {
			continue
		}

		go irc_clients.HandleNewIrcClient(conn)
	}
}
