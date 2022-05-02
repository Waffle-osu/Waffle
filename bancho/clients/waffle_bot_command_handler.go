package clients

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/packets"
	"strings"
)

type WaffleCommand interface {
	Execute(args []string) []string
}

var commandHandlers map[string]func(client_manager.OsuClient, []string) []string

func WaffleBotInitializeCommands() {
	commandHandlers = make(map[string]func(sender client_manager.OsuClient, args []string) []string)

	commandHandlers["!help"] = WaffleBotCommandHelp
	commandHandlers["!announce"] = WaffleBotCommandAnnounce
	commandHandlers["!roll"] = WaffleBotCommandRoll
}

func (client *Client) WaffleBotHandleCommand(sender client_manager.OsuClient, message packets.Message) {
	publicCommand := message.Target[0] == '#'

	var command string
	var arguments []string

	splitMessage := strings.Split(message.Message, " ")

	if len(splitMessage) == 0 {
		return
	}

	command = splitMessage[0]
	arguments = splitMessage[1:]

	handler, exists := commandHandlers[command]

	if exists == false {
		return
	}

	result := handler(sender, arguments)

	for _, messageString := range result {
		if publicCommand {
			channel, exists := chat.GetChannelByName(message.Target)

			if exists {
				channel.SendMessage(client, messageString, message.Target)
			}
		} else {
			packets.BanchoSendIrcMessage(sender.GetPacketQueue(), packets.Message{
				Sender:  "WaffleBot",
				Message: messageString,
				Target:  message.Target,
			})
		}
	}
}
