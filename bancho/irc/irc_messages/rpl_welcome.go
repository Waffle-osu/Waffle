package irc_messages

func IrcSendWelcome(message string) Message {
	return Message{
		NumCommand: 1,
		Trailing:   message,
	}
}
