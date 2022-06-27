package irc_clients

import (
	"Waffle/database"
	"bufio"
	"net"
)

var MOTD string = "" +
	" _       __      __________     __" +
	"| |     / /___ _/ __/ __/ /__  / /" +
	"| | /| / / __ `/ /_/ /_/ / _ \\/ / " +
	"| |/ |/ / /_/ / __/ __/ /  __/_/  " +
	"|__/|__/\\__,_/_/ /_/ /_/\\___(_)   " +
	"                                 "

type IrcClient struct {
	connection      net.Conn
	reader          *bufio.Reader
	continueRunning bool

	//Name used to address you on IRC
	//Must be unique across the network
	//This is the username used in /kick commands and similar
	Nickname string

	//This is used to populate the real name field when using /whois
	//can contain most characters
	Realname string

	//Is mainly used for people using 1 computer for more than 1 IRC User
	//To differenciate between them.
	//Also cannot be changed without a reconnect
	Username string

	//Password provided by the IRC Client
	Password string

	UserData database.User
}
