package irc_clients

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/irc/irc_messages"
	"strings"
)

func (client *IrcClient) SendChannelNames(channel *chat.Channel) {
	channelPrefix := "="

	if channel.ReadPrivileges > chat.PrivilegesNormal {
		channelPrefix = "@"
	}

	nameLines := []string{}

	nameString := ""

	for _, client := range channel.Clients {
		userPrefix := ""

		if client.GetUserPrivileges() > chat.PrivilegesNormal {
			userPrefix = "@"
		}

		nameString += userPrefix + strings.ReplaceAll(client.GetUsername(), " ", "_") + " "

		if len(nameString) > 255 {
			nameLines = append(nameLines, nameString)
			nameString = ""
		}
	}

	if len(nameString) != 0 {
		nameLines = append(nameLines, nameString)
	}

	for _, line := range nameLines {
		client.packetQueue <- irc_messages.IrcSendNameReply(channelPrefix, channel.Name, line)
	}

	client.packetQueue <- irc_messages.IrcSendEndOfNames(channel.Name, "End of NAMES")
}

func (client *IrcClient) SendWhoIs(checkClient client_manager.WaffleClient) {
	lastRecieve, logonTime := checkClient.GetIdleTimes()

	client.packetQueue <- irc_messages.IrcSendWhoIsUser(checkClient.GetUserData().Username)
	client.packetQueue <- irc_messages.IrcSendWhoIsChannels(checkClient.GetUserData().Username, checkClient.joinedChannels)
	client.packetQueue <- irc_messages.IrcSendWhoIsIdle(checkClient.GetUserData().Username, lastRecieve, logonTime)
	client.packetQueue <- irc_messages.IrcSendWhoIsServer(checkClient.GetUserData().Username)

	if checkClient.GetUserData().Privileges > 0 {
		client.packetQueue <- irc_messages.IrcSendWhoIsOperator(checkClient.GetUserData().Username)
	}

	client.packetQueue <- irc_messages.IrcSendEndOfWhoIs(checkClient.GetUserData().Username)
}
