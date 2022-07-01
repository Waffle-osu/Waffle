package irc_clients

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/irc/irc_messages"
	"Waffle/database"
	"Waffle/helpers"
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"net"
	"strings"
	"sync"
	"time"
)

func HandleNewIrcClient(connection net.Conn) {
	textReader := bufio.NewReader(connection)

	ircClient := IrcClient{
		connection:     connection,
		reader:         textReader,
		packetQueue:    make(chan irc_messages.Message, 128),
		joinedChannels: make(map[string]*chat.Channel),
		cleanMutex:     sync.Mutex{},
	}

	for ircClient.Username == "" || ircClient.Password == "" {
		line, err := textReader.ReadString('\n')

		if err != nil {
			return
		}

		helpers.Logger.Printf("[IRC@Debug] %s", line)

		message := irc_messages.ParseMessage(line)

		if len(message.Source) != 0 {
			helpers.Logger.Printf("[IRC@Debug] -- Source: %s", message.Source)
		}

		helpers.Logger.Printf("[IRC@Debug] -- Command: %s", message.Command)
		helpers.Logger.Printf("[IRC@Debug] -- Params: %s", strings.Join(message.Params, ", "))

		if len(message.Trailing) != 0 {
			helpers.Logger.Printf("[IRC@Debug] -- Trailing: %s", message.Trailing)
		}

		ircClient.ProcessMessage(message)
	}

	//TODO: irc tokens

	passwordHashed := md5.Sum([]byte(ircClient.Password))
	passwordHashedString := hex.EncodeToString(passwordHashed[:])

	userId, authResult := database.AuthenticateUser(ircClient.Username, passwordHashedString)

	if !authResult {
		ircClient.packetQueue <- irc_messages.IrcSendPasswordMismatch("Invalid Login!")

		ircClient.SendOffMessagesAndClose()
		return
	}

	queryResult, foundUser := database.UserFromDatabaseById(uint64(userId))

	if queryResult == -1 {
		ircClient.packetQueue <- irc_messages.IrcSendPasswordMismatch("Invalid Login!")

		ircClient.SendOffMessagesAndClose()
		return
	}

	if queryResult == -2 {
		ircClient.packetQueue <- irc_messages.IrcSendPasswordMismatch("Server Error.")

		ircClient.SendOffMessagesAndClose()
		return
	}

	ircClient.UserData = foundUser

	if ircClient.UserData.Banned == 1 {
		ircClient.packetQueue <- irc_messages.IrcSendPasswordMismatch("Login Error. Banned.")
		ircClient.packetQueue <- irc_messages.IrcSendBannedFromChan("You're banned!", "#osu")

		ircClient.SendOffMessagesAndClose()
		return
	}

	foundUsernameClient := client_manager.GetClientByName(ircClient.Username)

	if foundUsernameClient != nil {
		ircClient.packetQueue <- irc_messages.IrcSendPasswordMismatch("Login Error. Duplicate Usernames")
		ircClient.packetQueue <- irc_messages.IrcSendNicknameInUse(ircClient.Username, "Nickname already registered on server!")

		ircClient.SendOffMessagesAndClose()
		return
	}

	ircClient.packetQueue <- irc_messages.IrcSendTopic("#osu", "beyley is cute")
	ircClient.packetQueue <- irc_messages.IrcSendMotdBegin()

	for _, value := range MOTD {
		ircClient.packetQueue <- irc_messages.IrcSendMotd(value)
	}

	ircClient.packetQueue <- irc_messages.IrcSendMotdEnd()

	ircClient.lastPing = time.Now()
	ircClient.lastReceive = time.Now()
	ircClient.continueRunning = true

	go ircClient.HandleIncoming()
	go ircClient.SendOutgoing()
	go ircClient.MaintainClient()
}

func (client *IrcClient) SendOffMessagesAndClose() {
	for len(client.packetQueue) != 0 {
		formatted, _ := (<-client.packetQueue).FormatMessage(client.Username)

		client.connection.Write([]byte(formatted))
	}

	client.connection.Close()
}
