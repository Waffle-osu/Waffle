package irc_messages

import "strings"

func IrcSendPart(username string, channelName string) Message {
	return Message{
		Source:  strings.ReplaceAll(username, " ", "_") + "!" + username,
		Command: "PART",
		Params: []string{
			channelName,
		},
		SkipUsername: true,
	}
}
