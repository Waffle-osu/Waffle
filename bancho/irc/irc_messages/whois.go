package irc_messages

import (
	"strconv"
	"time"
)

func IrcSendWhoIsUser(username string) Message {
	return Message{
		NumCommand: RplWhoIsUser,
		Params: []string{
			username,
			username,
			"irc.waffle.nya",
			"*",
		},
		Trailing: username,
	}
}

func IrcSendWhoIsServer(username string) Message {
	return Message{
		NumCommand: RplWhoIsServer,
		Params: []string{
			username,
			"irc.waffle.nya",
		},
		Trailing: "Chillin' on Waffle! https://github.com/Eeveelution/Waffle",
	}
}

func IrcSendWhoIsOperator(username string) Message {
	return Message{
		NumCommand: RplWhoIsOperator,
		Params: []string{
			username,
		},
		Trailing: "is a Moderator",
	}
}

func IrcSendWhoIsChannels(username string, channelString string) Message {
	return Message{
		NumCommand: RplWhoIsChannels,
		Params: []string{
			username,
		},
		Trailing: channelString,
	}
}

func IrcSendWhoIsIdle(username string, lastRecieve time.Time, signon time.Time) Message {
	lastSeconds := int64(time.Since(lastRecieve).Seconds())
	signOnUnix := signon.Unix()

	return Message{
		NumCommand: RplWhoIsIdle,
		Params: []string{
			username,
			strconv.FormatInt(lastSeconds, 10),
			strconv.FormatInt(signOnUnix, 10),
		},
		Trailing: "Seconds since last Recieve, Unix time of sign-in",
	}
}

func IrcSendEndOfWhoIs(username string) Message {
	return Message{
		NumCommand: RplEndOfWhoIs,
		Params: []string{
			username,
		},
		Trailing: "End of WHOIS list.",
	}
}
