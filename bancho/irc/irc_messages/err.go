package irc_messages

func IrcSendPasswordMismatch(message string) Message {
	return Message{
		NumCommand: ErrPasswdMissmatch,
		Trailing:   message,
	}
}

func IrcSendNoSuchChannel(message string, channel string) Message {
	return Message{
		NumCommand: ErrBannedFromChan,
		Params: []string{
			channel,
		},
		Trailing: message,
	}
}

func IrcSendNicknameInUse(username string, message string) Message {
	return Message{
		NumCommand: ErrNicknameInUse,
		Params: []string{
			username,
		},
		Trailing: message,
	}
}

func IrcSendNoSuchNick(username string) Message {
	return Message{
		NumCommand: ErrNoSuchNick,
		Params: []string{
			username,
		},
		Trailing: "No such nickname/channel!",
	}
}

func IrcSendBannedFromChan(message string, channel string) Message {
	return Message{
		NumCommand: ErrBannedFromChan,
		Params: []string{
			channel,
		},
		Trailing: message,
	}
}

func IrcSendAlreadyRegistered(message string) Message {
	return Message{
		NumCommand: ErrAlreadyRegistered,
		Trailing:   message,
	}
}

func IrcSendNoSuchServer(server string) Message {
	return Message{
		NumCommand: ErrNoSuchServer,
		Params: []string{
			server,
		},
		Trailing: "No such server!",
	}
}
