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

		time.Sleep(time.Second)
	}
}

func (client *Client) WaffleBotHandleOutgoing() {
	for packet := range client.PacketQueue {
		packetDataReader := bytes.NewBuffer(packet.PacketData)

		switch packet.PacketId {
		case packets.BanchoSendMessage:
			message := packets.ReadMessage(packetDataReader)

			if message.Message[0] == '!' {
				sender := client_manager.GetClientByName(message.Sender)
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
