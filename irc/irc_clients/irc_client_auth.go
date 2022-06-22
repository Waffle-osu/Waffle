package irc_clients

import (
	"Waffle/helpers"
	"Waffle/irc/irc_reading"
	"bufio"
	"net"
	"strings"
)

func HandleNewIrcClient(connection net.Conn) {
	textReader := bufio.NewReader(connection)

	ircClient := IrcClient{
		Connection: connection,
		Reader:     textReader,
	}

	for i := 0; i != 16; i++ {
		line, err := textReader.ReadString('\n')

		if err != nil {
			return
		}

		helpers.Logger.Printf("[IRC@Debug] %s", line)

		message := irc_reading.ParseMessage(line)

		if len(message.Source) != 0 {
			helpers.Logger.Printf("[IRC@Debug] -- Source: %s", message.Source)
		}

		helpers.Logger.Printf("[IRC@Debug] -- Command: %s", message.Command)
		helpers.Logger.Printf("[IRC@Debug] -- Params: %s", strings.Join(message.Params, ", "))

		if len(message.Trailing) != 0 {
			helpers.Logger.Printf("[IRC@Debug] -- Trailing: %s", message.Trailing)
		}
	}
}
