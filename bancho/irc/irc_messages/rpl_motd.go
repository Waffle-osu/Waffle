package irc_messages

func IrcSendMotd(motd string) Message {
	return Message{
		NumCommand: RplMotd,
		Trailing:   "- " + motd,
	}
}
