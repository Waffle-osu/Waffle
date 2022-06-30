package irc_clients

import (
	"Waffle/bancho/irc/irc_messages"
	"strings"
)

func (client *IrcClient) ProcessMessage(message irc_messages.Message) {
	switch message.Command {
	case "NICK":
		client.Nickname = strings.Join(message.Params, " ")
	case "USER":
		client.Username = message.Params[0]
		client.Realname = message.Trailing
	case "PASS":
		if message.Trailing == "" {
			client.Password = strings.Join(message.Params, " ")
		} else {
			client.Password = message.Trailing
		}
	}
}

func (client *IrcClient) HandleIncoming() {
	for client.continueRunning {

	}
}
