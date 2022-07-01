package irc_messages

import "strings"

func IrcSendPrivMsg(username string, target string, message string) Message {
	return Message{
		SkipUsername: true,
		Command:      "PRIVMSG",
		Source:       strings.ReplaceAll(username, " ", "_") + "!" + username,
		Params: []string{
			target,
		},
		Trailing: message,
	}
}
