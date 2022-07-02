package irc_messages

func IrcSendEndOfWho(query string) Message {
	return Message{
		NumCommand: RplEndOfWho,
		Params: []string{
			query,
		},
		Trailing: "End of WHO list",
	}
}