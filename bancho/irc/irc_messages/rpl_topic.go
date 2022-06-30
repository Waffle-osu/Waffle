package irc_messages

func IrcSendTopic(channel string, topic string) Message {
	return Message{
		NumCommand: RplTopic,
		Params: []string{
			channel,
		},
		Trailing: topic,
	}
}
