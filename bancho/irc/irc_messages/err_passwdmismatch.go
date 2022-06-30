package irc_messages

func IrcSendPasswordMismatch(message string) Message {
	return Message{
		NumCommand: ErrPasswdMissmatch,
		Trailing:   message,
	}
}
