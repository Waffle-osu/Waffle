package irc

import (
	"Waffle/config"
	"Waffle/helpers"
	"Waffle/irc/irc_clients"
	"crypto/tls"
)

func RunIrcSSL() {
	if config.HostIrcSsl == "false" || config.SSLCertLocation == "" || config.SSLKeyLocation == "" {
		return
	}

	cert, certErr := tls.LoadX509KeyPair(config.SSLCertLocation, config.SSLKeyLocation)

	if certErr != nil {
		helpers.Logger.Fatalf("[IRC/SSL] Failed to create certificate.")
	}

	certConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	listener, err := tls.Listen("tcp", config.IrcSslIp, certConfig)

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

		go irc_clients.HandleNewIrcClient(conn)
	}
}
