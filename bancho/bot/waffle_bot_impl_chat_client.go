package bot

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/osu/base_packet_structures"
)

func (waffleBot *WaffleBot) GetUserPrivileges() int32 {
	return waffleBot.UserData.Privileges
}

func (waffleBot *WaffleBot) GetUsername() string {
	return "WaffleBot"
}

func (waffleBot *WaffleBot) InformChannelJoin(chatClient chat.ChatClient, channel *chat.Channel) {

}

func (waffleBot *WaffleBot) InformChannelPart(chatClient chat.ChatClient, channel *chat.Channel) {

}

func (waffleBot *WaffleBot) SendChatMessage(sender string, content string, channel string) {
	waffleBot.BanchoIrcMessage(base_packet_structures.Message{
		Message: content,
		Target:  channel,
		Sender:  sender,
	})
}

func (waffleBot *WaffleBot) GetSilencedUntilUnix() int64 {
	return 0
}

func (client *WaffleBot) SetSilencedUntilUnix(untilUnix int64) {

}
