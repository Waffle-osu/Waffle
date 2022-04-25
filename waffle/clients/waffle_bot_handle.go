package clients

import (
	"Waffle/waffle/chat"
	"Waffle/waffle/client_manager"
	"Waffle/waffle/packets"
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

func (client *Client) WaffleBotMaintainClient() {
	for client.continueRunning {
		//Maybe i'll add some fancy stuff here like funny statuses but as it stands this will be empty

		time.Sleep(time.Second)
	}
}

// WaffleBotHandleOutgoing Handles stuff that's been sent to WaffleBot
func (client *Client) WaffleBotHandleOutgoing() {
	for packet := range client.PacketQueue {
		packetDataReader := bytes.NewBuffer(packet.PacketData)

		switch packet.PacketId {
		case packets.BanchoSendMessage:
			message := packets.ReadMessage(packetDataReader)

			//Handles commands
			if message.Message[0] == '!' {
				sender := client_manager.GetClientByName(message.Sender)
				//This determines whether the response to the command will be sent publicly in chat or privately in DMs
				publicCommand := message.Target[0] == '#'

				splitCommand := strings.Split(message.Message, " ")

				if len(splitCommand) == 0 {
					break
				}

				switch strings.ToLower(splitCommand[0]) {
				case "!help":
					if publicCommand {
						channel, exists := chat.GetChannelByName(message.Target)

						if exists {
							channel.SendMessage(client, "Currently, there's not much to talk about as this is in very early stages... expect there to come more as this develops", message.Target)
						}
					} else {
						packets.BanchoSendIrcMessage(sender.GetPacketQueue(), packets.Message{
							Sender:  "WaffleBot",
							Message: "Currently, there's not much to talk about as this is in very early stages... expect there to come more as this develops",
							Target:  message.Target,
						})
					}
				}
			}
			break
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
			break
		default:
			fmt.Printf("WaffleBot got %s\n", packets.GetPacketName(packet.PacketId))
			break
		}
	}
}
