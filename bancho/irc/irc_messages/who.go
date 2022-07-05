package irc_messages

func IrcSendWhoReply(query string, username string, isAway bool, privileges int32) Message {
	flags := ""

	if isAway {
		flags += "H"
	} else {
		flags += "G"
	}

	if privileges >= 1 {
		flags += "*"
	}

	return Message{
		NumCommand: RplWhoReply,
		Params: []string{
			query,
			username,
			"irc/waffle/nya",
			"irc.waffle.nya",
			username,
			flags,
		},
		Trailing: "0 " + username,
	}
}

func IrcSendEndOfWho(query string) Message {
	return Message{
		NumCommand: RplEndOfWho,
		Params: []string{
			query,
		},
		Trailing: "End of WHO list",
	}
}
