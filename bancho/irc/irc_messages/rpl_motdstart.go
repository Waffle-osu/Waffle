package irc_messages

func IrcSendMotdBegin() Message {
	return Message{
		NumCommand: RplMotdStart,
		Trailing: "-",
	}
}