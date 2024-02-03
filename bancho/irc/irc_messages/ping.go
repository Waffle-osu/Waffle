package irc_messages

func IrcSendPing(token string, skipSource bool) Message {
	return Message{
		Command: "PING",
		Params: []string{
			token,
		},
		SkipSource: skipSource,
	}
}

func IrcSendPong(token string) Message {
	return Message{
		Command: "PONG",
		Params: []string{
			"irc.waffle.nya",
			token,
		},
	}
}
