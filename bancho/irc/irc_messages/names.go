package irc_messages

func IrcSendNameReply(channelPrefix string, channelName string, names string) Message {
	return Message{
		NumCommand: RplNameReply,
		Params: []string{
			channelPrefix,
			channelName,
		},
		Trailing: names,
	}
}

func IrcSendEndOfNames(channel string, message string) Message {
	return Message{
		NumCommand: RplEndOfNames,
		Params: []string{
			channel,
		},
		Trailing: message,
	}
}
