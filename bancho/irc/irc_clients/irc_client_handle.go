package irc_clients

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/clients"
	"Waffle/bancho/irc/irc_messages"
	"Waffle/bancho/lobby"
	"Waffle/bancho/misc"
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/helpers"
	"context"
	"fmt"
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
		channels := []string{}

		for _, requestedChannel := range message.Params {
			channels = append(channels, strings.Split(requestedChannel, ",")...)
		}

		for _, channel := range channels {
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
	case "PART":
		channels := []string{}

		for _, requestedChannel := range message.Params {
			channels = append(channels, strings.Split(requestedChannel, ",")...)
		}

		for _, channel := range channels {
			joinedChannel, exists := client.joinedChannels[channel]

			if exists {
				joinedChannel.Leave(client)
			} else {
				client.packetQueue <- irc_messages.IrcSendNotOnChannel(channel)
			}
		}
	case "NAMES":
		channels := []string{}

		for _, requestedChannel := range message.Params {
			channels = append(channels, strings.Split(requestedChannel, ",")...)
		}

		for _, channel := range channels {
			foundChannel, exists := chat.GetChannelByName(channel)

			if exists {
				client.SendChannelNames(foundChannel)
			} else {
				client.packetQueue <- irc_messages.IrcSendNoSuchChannel("No such channel!", channel)
			}
		}
	case "QUIT":
		client.CleanupClient(message.Trailing)

		client_manager.BroadcastPacketOsu(func(_client client_manager.WaffleClient) {
			_client.BanchoHandleIrcQuit(client.Username)
		})
	case "PRIVMSG":
		if len(message.Params) != 0 {
			if time.Now().Unix() < int64(client.UserData.SilencedUntil) {
				client.SendChatMessage("WaffleBot", fmt.Sprintf("You're silenced for at least %d seconds!", int64(client.UserData.SilencedUntil)-time.Now().Unix()), client.UserData.Username)
			} else {
				messageText := message.Trailing
				target := message.Params[0]

				//Chhannel into which the message gets sent
				//WaffleBot/Lobby commands get sent there awell
				//Aswell as their responses
				var sendChannel *chat.Channel

				//Command outputs
				var returnMessages []string

				//Commands start with !
				if messageText[0] == '!' {
					//MP commands take a different route
					//because the handler has to take in a LobbyClient
					//instead of just a WaffleClient
					if strings.HasPrefix(messageText, "!mp") {
						returnMessages = lobby.LobbyHandleCommandMultiplayer(client, messageText)
					} else {
						returnMessages = clients.WaffleBotInstance.WaffleBotHandleCommand(client, base_packet_structures.Message{
							Sender:  client.Username,
							Message: messageText,
							Target:  target,
						})
					}
				}

				//Reroute for multiplayer, or find channel
				if target == "#multiplayer" {
					if client.currentMultiLobby != nil {
						sendChannel = client.currentMultiLobby.MultiChannel
					} else {
						foundChannel, exists := client.joinedChannels[target]

						if exists {
							sendChannel = foundChannel
						} else {
							sendChannel = nil
						}
					}
				}

				if sendChannel != nil {
					sendChannel.SendMessage(client, messageText, target)

					for _, content := range returnMessages {
						sendChannel.SendMessage(clients.WaffleBotInstance, content, target)
					}
				} else {
					//Could also be a client.
					//So before sending a error message, we need to search for a client
					foundClient := client_manager.GetClientByName(target)

					if foundClient != nil {
						foundClient.BanchoIrcMessage(base_packet_structures.Message{
							Sender:  client.Username,
							Target:  target,
							Message: messageText,
						})

						for _, content := range returnMessages {
							foundClient.BanchoIrcMessage(base_packet_structures.Message{
								Sender:  "WaffleBot",
								Target:  target,
								Message: content,
							})
						}

						break
					} else {
						client.packetQueue <- irc_messages.IrcSendNoSuchChannel("Channel either doesn't exist or you haven't joined it. No user under such Username could be found either.", message.Params[0])

						break
					}
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
	case "LIST":
		client.packetQueue <- irc_messages.IrcSendListStart()

		if len(message.Params) == 0 {
			for _, channel := range chat.GetAvailableChannels() {
				if (channel.ReadPrivileges & client.GetUserPrivileges()) <= 0 {
					continue
				}

				client.packetQueue <- irc_messages.IrcSendListReply(channel)
			}
		} else {
			joinedQuery := strings.Join(message.Params, " ")

			if strings.Contains(joinedQuery, "#") {
				requestedChannels := strings.Split(joinedQuery, ",")

				for _, requestedChannel := range requestedChannels {
					foundChannel, exists := chat.GetChannelByName(requestedChannel)

					if (foundChannel.ReadPrivileges & client.GetUserPrivileges()) <= 0 {
						continue
					}

					if exists {
						client.packetQueue <- irc_messages.IrcSendListReply(foundChannel)
					}
				}
			} else {
				for _, channel := range chat.GetAvailableChannels() {
					if (channel.ReadPrivileges & client.GetUserPrivileges()) <= 0 {
						continue
					}

					client.packetQueue <- irc_messages.IrcSendListReply(channel)
				}
			}
		}

		client.packetQueue <- irc_messages.IrcSendListEnd()
	case "PING":
		token := message.Params[0]

		client.packetQueue <- irc_messages.IrcSendPong(token)
	case "PONG":
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

func (client *IrcClient) MaintainClient(ctx context.Context) {
	pingTicker := time.NewTicker(PingTimeout * time.Second)
	receiveTicker := time.NewTicker(ReceiveTimeout * time.Second)

	for {
		select {
		case <-ctx.Done():
			//We close in MaintainClient instead of in CleanupClient to avoid possible double closes, causing panics
			helpers.Logger.Printf("[IRC@Handling] Closed %s's Packet Queue", client.UserData.Username)

			close(client.packetQueue)

			pingTicker.Stop()
			receiveTicker.Stop()
			return
		case <-pingTicker.C:
			client.lastPingToken = fmt.Sprintf("irc.waffle.nya@%d", time.Now().Unix())

			client.packetQueue <- irc_messages.IrcSendPing(client.lastPingToken)

			client.lastPing = time.Now()
		case <-receiveTicker.C:
			lastReceiveUnix := client.lastReceive.Unix()
			unixNow := time.Now().Unix()

			if lastReceiveUnix+ReceiveTimeout <= unixNow {
				client.CleanupClient("Client Timed out.")
			}
		case message := <-client.packetQueue:
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
}
