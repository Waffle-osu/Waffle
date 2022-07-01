package irc_messages

import "strings"

func IrcSendJoin(username string, channelName string) Message {
	return Message{
		Source:       strings.ReplaceAll(username, " ", "_") + "!" + username,
		Command:      "JOIN",
		Trailing:     channelName,
		SkipUsername: true,
	}
}
