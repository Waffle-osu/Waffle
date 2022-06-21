package irc_clients

import (
	"bufio"
	"net"
)

func HandleNewIrcSslClient(connection net.Conn) {
	textReader := bufio.NewReader(connection)

	textReader.ReadString()
}
