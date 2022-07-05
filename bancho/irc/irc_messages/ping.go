package irc_messages

func IrcSendPing(token string) Message {
	return Message{
		Command: "PING",
		Params: []string{
			token,
		},
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
