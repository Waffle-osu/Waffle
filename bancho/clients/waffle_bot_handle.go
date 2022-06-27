package clients

import (
	"time"
)

func (client *WaffleBot) WaffleBotMaintainClient() {
	//for client.continueRunning {
	//Maybe i'll add some fancy stuff here like funny statuses but as it stands this will be empty

	time.Sleep(time.Second)
	//}
}

// WaffleBotHandleOutgoing Handles stuff that's been sent to WaffleBot
func (client *WaffleBot) WaffleBotHandleOutgoing() {
	/*
		for packet := range client.PacketQueue {
			packetDataReader := bytes.NewBuffer(packet.PacketData)

			switch packet.PacketId {
			case packets.BanchoSendMessage:
				message := packets.ReadMessage(packetDataReader)
				//Handles commands
				if message.Message[0] == '!' {
					sender := client_manager.GetClientByName(message.Sender)

					go client.WaffleBotHandleCommand(sender, message)
				}
			case packets.OsuSendIrcMessagePrivate:
				message := packets.ReadMessage(packetDataReader)
				//Assign a sender, as the client doesn't seem to send itself as the sender
				message.Sender = client.UserData.Username
				//Handles commands
				if message.Message[0] == '!' {
					sender := client_manager.GetClientByName(message.Sender)

					go client.WaffleBotHandleCommand(sender, message)
				}
			case packets.BanchoSpectatorJoined:
				var userId int32

				binary.Read(packetDataReader, binary.LittleEndian, &userId)

				userById := client_manager.GetClientById(userId)

				//Fun little easter egg, just thought i'd add it in initially for testing if it can recieve packets
				if userById != nil {
					packets.BanchoSendIrcMessage(userById.GetPacketQueue(), packets.Message{
						Sender:  "WaffleBot",
						Message: "Currently, there's not point in spectating me!!",
						Target:  userById.GetUserData().Username,
					})
				}
			default:
				helpers.Logger.Printf("[Bancho@WaffleBotHandle] WaffleBot got %s\n", packets.GetPacketName(packet.PacketId))
			}
		}
	*/

}
