package irc_messages

func IrcSendBannedFromChan(message string, channel string) Message {
	return Message{
		NumCommand: ErrBannedFromChan,
		Params: []string{
			channel,
		},
		Trailing: message,
	}
}
