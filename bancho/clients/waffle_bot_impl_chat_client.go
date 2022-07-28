package clients

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/common"
)

func (waffleBot *WaffleBot) GetClientType() common.ClientType {
	return common.ClientTypeOsu1816
}

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
