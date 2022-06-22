package irc

import (
	"Waffle/config"
	"Waffle/helpers"
	"net"
)

func RunIrcSSL() {
	if config.HostIrcSsl == "false" || config.SSLCertLocation == "" || config.SSLKeyLocation == "" {
		return
	}

	listener, err := net.Listen("tcp", config.IrcSslIp)

	if err != nil {
		helpers.Logger.Printf("[IRC/SSL] Failed to create TCP Listener for IRC/SSL on %s\n", config.IrcSslIp)
	}

	helpers.Logger.Printf("Running IRC/SSL on %s\n", config.IrcSslIp)

	for {
		conn, err := listener.Accept()

		helpers.Logger.Printf("[IRC/SSL] Accepted Connection!\n")

		if err != nil {
			continue
		}

	}
}
