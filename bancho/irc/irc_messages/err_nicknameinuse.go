package irc_messages

func IrcSendNicknameInUse(username string, message string) Message {
	return Message{
		NumCommand: ErrNicknameInUse,
		Params: []string{
			username,
		},
		Trailing: message,
	}
}
