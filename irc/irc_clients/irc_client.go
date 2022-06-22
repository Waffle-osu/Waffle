package irc_clients

import (
	"bufio"
	"net"
)

type IrcClient struct {
	connection      net.Conn
	reader          *bufio.Reader
	continueRunning bool

	//Name used to address you on IRC
	//Must be unique across the network
	Nickname string
	Realname string
	Username string
}
