package irc_messages

func IrcSendAlreadyRegistered(message string) Message {
	return Message{
		NumCommand: ErrAlreadyRegistered,
		Trailing:   message,
	}
}
