package irc_messages

func IrcSendNoSuchChannel(message string, channel string) Message {
	return Message{
		NumCommand: ErrBannedFromChan,
		Params: []string{
			channel,
		},
		Trailing: message,
	}
}
