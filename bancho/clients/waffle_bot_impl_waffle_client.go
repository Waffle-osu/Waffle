package clients

import (
	"Waffle/bancho/client_manager"
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/database"
	"time"
)

func (waffleBot *WaffleBot) GetUserId() int32 {
	return int32(waffleBot.UserData.UserID)
}

func (waffleBot *WaffleBot) GetRelevantUserStats() database.UserStats {
	return waffleBot.OsuStats
}

func (waffleBot *WaffleBot) GetUserStatus() base_packet_structures.StatusUpdate {
	return waffleBot.Status
}

func (waffleBot *WaffleBot) GetUserData() database.User {
	return waffleBot.UserData
}

func (waffleBot *WaffleBot) GetClientTimezone() int32 {
	return 0
}

func (waffleBot *WaffleBot) GetIdleTimes() (lastReceive time.Time, logonTime time.Time) {
	return waffleBot.lastReceive, waffleBot.logonTime
}

func (waffleBot *WaffleBot) GetFormattedJoinedChannels() string {
	channelString := ""

	for _, value := range waffleBot.joinedChannels {
		if value.ReadPrivileges == 0 {
			channelString += value.Name + " "
		}
	}

	return channelString
}

func (waffleBot *WaffleBot) CleanupClient(reason string) {

}

func (waffleBot *WaffleBot) Cut() {

}

func (waffleBot *WaffleBot) GetAwayMessage() string {
	return waffleBot.awayMessage
}

func (waffleBot *WaffleBot) BanchoHandleOsuQuit(userId int32, username string) {

}

func (waffleBot *WaffleBot) BanchoHandleIrcQuit(username string) {

}

func (waffleBot *WaffleBot) BanchoSpectatorJoined(userId int32) {
	userById := client_manager.GetClientById(userId)

	//Fun little easter egg, just thought i'd add it in initially for testing if it can recieve packets
	if userById != nil {
		userById.BanchoIrcMessage(base_packet_structures.Message{
			Sender:  "WaffleBot",
			Message: "Currently, there's not point in spectating me!!",
			Target:  userById.GetUserData().Username,
		})
	}
}

func (waffleBot *WaffleBot) BanchoSpectatorLeft(userId int32) {

}

func (waffleBot *WaffleBot) BanchoFellowSpectatorJoined(userId int32) {

}

func (waffleBot *WaffleBot) BanchoFellowSpectatorLeft(userId int32) {

}

func (waffleBot *WaffleBot) BanchoSpectatorCantSpectate(userId int32) {

}

func (waffleBot *WaffleBot) BanchoSpectateFrames(frameBundle base_packet_structures.SpectatorFrameBundle) {

}

func (waffleBot *WaffleBot) BanchoIrcMessage(message base_packet_structures.Message) {
	if message.Target == "WaffleBot" {
		if message.Message[0] == '!' {
			client := client_manager.GetClientByName(message.Sender)

			if client == nil {
				return
			}

			returnMessages := waffleBot.WaffleBotHandleCommand(client, message)

			for _, content := range returnMessages {
				client.BanchoIrcMessage(base_packet_structures.Message{
					Sender:  "WaffleBot",
					Message: content,
					Target:  "WaffleBot",
				})
			}
		}
	}
}

func (waffleBot *WaffleBot) BanchoOsuUpdate(stats database.UserStats, update base_packet_structures.StatusUpdate) {

}

func (waffleBot *WaffleBot) BanchoPresence(user database.User, stats database.UserStats, timezone int32) {

}

func (waffleBot *WaffleBot) BanchoAnnounce(message string) {

}

func (waffleBot *WaffleBot) BanchoGetAttention() {

}
