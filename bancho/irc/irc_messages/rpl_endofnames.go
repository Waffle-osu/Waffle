package irc_messages

func IrcSendEndOfNames(channel string, message string) Message {
	return Message{
		NumCommand: RplEndOfNames,
		Params: []string{
			channel,
		},
		Trailing: message,
	}
}
