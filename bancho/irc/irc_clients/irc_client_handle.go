package irc_clients

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/irc/irc_messages"
	"Waffle/bancho/misc"
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/helpers"
	"strings"
	"time"
)

func (client *IrcClient) ProcessMessage(message irc_messages.Message, rawLine string) {
	switch strings.ToUpper(message.Command) {
	case "NICK":
		client.Nickname = strings.Join(message.Params, " ")

		//TODO: BanchoHandleIrcChangeUsername
	case "USER":
		if client.Username == "" && client.Realname == "" {
			client.Username = message.Params[0]
			client.Realname = message.Trailing
		} else {
			client.packetQueue <- irc_messages.IrcSendAlreadyRegistered("You may not reregister")
		}
	case "PASS":
		if client.Password == "" {
			if message.Trailing == "" {
				client.Password = strings.Join(message.Params, " ")
			} else {
				client.Password = message.Trailing
			}
		} else {
			client.packetQueue <- irc_messages.IrcSendAlreadyRegistered("You may not reregister")
		}
	case "JOIN":
		for _, channel := range message.Params {
			foundChannel, exists := chat.GetChannelByName(channel)

			if !exists {
				client.packetQueue <- irc_messages.IrcSendNoSuchChannel("No such channel!", channel)
				return
			}

			success := foundChannel.Join(client)

			if success {
				client.joinedChannels[foundChannel.Name] = foundChannel

				client.packetQueue <- irc_messages.IrcSendTopic(channel, foundChannel.Description)

				client.SendChannelNames(foundChannel)
			} else {
				client.packetQueue <- irc_messages.IrcSendBannedFromChan("Joining channel failed.", channel)
			}
		}
	case "PRIVMSG":
		if len(message.Params) != 0 {
			foundChannel, exists := client.joinedChannels[message.Params[0]]

			if exists {
				foundChannel.SendMessage(client, message.Trailing, message.Params[0])
			} else {
				foundClient := client_manager.GetClientByName(message.Params[0])

				if foundClient != nil {
					foundClient.BanchoIrcMessage(base_packet_structures.Message{
						Sender:  client.Username,
						Target:  message.Params[0],
						Message: message.Trailing,
					})
				} else {
					client.packetQueue <- irc_messages.IrcSendNoSuchChannel("Channel either doesn't exist or you haven't joined it. No user under such Username could be found either.", message.Params[0])
				}
			}
		}
	case "WHO":
		query := message.Params[0]

		channelQuery := query[0] == '#'

		if channelQuery {
			foundChannel, exists := client.joinedChannels[query]

			if exists {
				for _, channelClient := range foundChannel.Clients {
					isAway := channelClient.GetAwayMessage() == ""

					client.packetQueue <- irc_messages.IrcSendWhoReply(query, channelClient.GetUsername(), isAway, channelClient.GetUserPrivileges())
				}
			} else {
				client.packetQueue <- irc_messages.IrcSendNoSuchChannel("Channel either doesn't exist or you haven't joined it. No user under such Username could be found either.", message.Params[0])
			}
		} else {
			var foundClient client_manager.WaffleClient = nil

			//Regular Username query
			if foundClient == nil {
				foundClient = client_manager.GetClientByName(query)
			}

			//Try again but replace _ by space
			if foundClient == nil {
				foundClient = client_manager.GetClientByName(strings.ReplaceAll(query, "_", " "))
			}

			if foundClient == nil {
				client.packetQueue <- irc_messages.IrcSendNoSuchChannel("Channel either doesn't exist or you haven't joined it. No user under such Username could be found either.", message.Params[0])
			} else {
				isAway := foundClient.GetAwayMessage() == ""

				client.packetQueue <- irc_messages.IrcSendWhoReply(query, foundClient.GetUserData().Username, isAway, foundClient.GetUserData().Privileges)
			}
		}

		client.packetQueue <- irc_messages.IrcSendEndOfWho(query)
	case "WHOIS":
		if len(message.Params) == 2 {
			server := message.Params[0]
			username := message.Params[1]

			if server != "irc.waffle.nya" {
				foundServerClient := client_manager.GetClientByName(server)

				if foundServerClient != nil {
					foundWhoIsClient := client_manager.GetClientByName(username)

					if foundWhoIsClient != nil {
						client.SendWhoIs(foundWhoIsClient)
					} else {
						client.packetQueue <- irc_messages.IrcSendNoSuchNick(username)
						client.packetQueue <- irc_messages.IrcSendEndOfWhoIs(username)
					}
				} else {
					client.packetQueue <- irc_messages.IrcSendNoSuchServer(server)
					client.packetQueue <- irc_messages.IrcSendEndOfWhoIs(username)
				}
			}

		} else if len(message.Params) == 1 {
			username := message.Params[0]

			foundWhoIsClient := client_manager.GetClientByName(username)

			if foundWhoIsClient != nil {
				client.SendWhoIs(foundWhoIsClient)
			} else {
				client.packetQueue <- irc_messages.IrcSendNoSuchNick(username)
				client.packetQueue <- irc_messages.IrcSendEndOfWhoIs(username)
			}
		}
	case "CAP":
	default:
		helpers.Logger.Printf("[IRC@Debug] UNHANDLED COMMAND: %s", rawLine)

		if len(message.Source) != 0 {
			helpers.Logger.Printf("[IRC@Debug] -- Source: %s", message.Source)
		}

		helpers.Logger.Printf("[IRC@Debug] -- Command: %s", message.Command)
		helpers.Logger.Printf("[IRC@Debug] -- Params: %s", strings.Join(message.Params, ", "))

		if len(message.Trailing) != 0 {
			helpers.Logger.Printf("[IRC@Debug] -- Trailing: %s", message.Trailing)
		}
	}
}

func (client *IrcClient) HandleIncoming() {
	for client.continueRunning {
		line, err := client.reader.ReadString('\n')

		if err != nil {
			return
		}

		client.lastReceive = time.Now()

		message := irc_messages.ParseMessage(line)

		client.ProcessMessage(message, line)
	}
}

func (client *IrcClient) SendOutgoing() {
	for message := range client.packetQueue {
		formatted, _ := message.FormatMessage(client.Username)

		asBytes := []byte(formatted)

		go func() {
			misc.StatsSendLock.Lock()
			misc.StatsBytesSent += uint64(len(asBytes))
			misc.StatsSendLock.Unlock()
		}()

		client.connection.Write(asBytes)
	}
}

func (client *IrcClient) MaintainClient() {
	for client.continueRunning {
		lastReceiveUnix := client.lastReceive.Unix()
		lastPingUnix := client.lastPing.Unix()
		unixNow := time.Now().Unix()

		if lastReceiveUnix+ReceiveTimeout <= unixNow {
			client.CleanupClient("Client Timed out.")

			client.continueRunning = false
		}

		if lastPingUnix+PingTimeout <= unixNow {
			//client.BanchoPing()

			client.lastPing = time.Now()
		}

		time.Sleep(time.Second)
	}

	//We close in MaintainClient instead of in CleanupClient to avoid possible double closes, causing panics
	helpers.Logger.Printf("[IRC@Handling] Closed %s's Packet Queue", client.UserData.Username)

	close(client.packetQueue)
}
