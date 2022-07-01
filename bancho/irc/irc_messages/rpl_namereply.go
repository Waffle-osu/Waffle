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
