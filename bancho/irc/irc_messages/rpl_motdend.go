package irc_messages

func IrcSendMotdEnd() Message {
	return Message{
		NumCommand: RplEndOfMotd,
		Trailing:   "-",
	}
}
