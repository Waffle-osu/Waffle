package irc_messages

import (
	"Waffle/bancho/chat"
	"fmt"
)

func IrcSendListStart() Message {
	return Message{
		NumCommand: RplListStart,
		Params: []string{
			"Channel",
		},
		Trailing: "Users  Name",
	}
}

func IrcSendListReply(channel *chat.Channel) Message {
	return Message{
		NumCommand: RplList,
		Params: []string{
			channel.Name,
			fmt.Sprintf("%d", len(channel.Clients)),
		},
		Trailing: channel.Description,
	}
}

func IrcSendListEnd() Message {
	return Message{
		NumCommand: RplListEnd,
		Params: []string{
			"Channel",
		},
		Trailing: "End of /LIST",
	}
}
