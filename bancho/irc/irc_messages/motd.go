package irc_messages

func IrcSendMotd(motd string) Message {
	return Message{
		NumCommand: RplMotd,
		Trailing:   "- " + motd,
	}
}

func IrcSendMotdEnd() Message {
	return Message{
		NumCommand: RplEndOfMotd,
		Trailing:   "-",
	}
}

func IrcSendMotdBegin() Message {
	return Message{
		NumCommand: RplMotdStart,
		Trailing: "-",
	}
}