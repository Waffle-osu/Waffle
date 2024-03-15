package irc_clients

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/irc/irc_messages"
	"Waffle/bancho/spectator"
	"Waffle/database"
	"bufio"
	"context"
	"crypto/md5"
	"encoding/hex"
	"net"
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
		spectators:     make(map[int32]spectator.SpectatorClient),
	}

	for ircClient.Username == "" || ircClient.Password == "" {
		line, err := textReader.ReadString('\n')

		if err != nil {
			return
		}

		message := irc_messages.ParseMessage(line)

		ircClient.ProcessMessage(message, line)
		ircClient.SendOffMessages()
	}

	//TODO: IRC OTP: https://github.com/Eeveelution/Waffle/issues/15

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

	foundUsernameClient := client_manager.ClientManager.GetClientByName(ircClient.Username)

	if foundUsernameClient != nil {
		ircClient.packetQueue <- irc_messages.IrcSendPasswordMismatch("Login Error. Duplicate Usernames")
		ircClient.packetQueue <- irc_messages.IrcSendNicknameInUse(ircClient.Username, "Nickname already registered on server!")

		ircClient.SendOffMessagesAndClose()
		return
	}

	if ircClient.IsOsu {
		statsQueryResult, stats := database.UserStatsFromDatabase(ircClient.UserData.UserID, 0)

		if statsQueryResult == 0 {
			ircClient.OsuStats = stats
		}
	}

	ircClient.packetQueue <- irc_messages.IrcSendTopic("#osu", "beyley is cute")
	ircClient.packetQueue <- irc_messages.IrcSendMotdBegin()

	for _, value := range MOTD {
		ircClient.packetQueue <- irc_messages.IrcSendMotd(value)
	}

	ircClient.packetQueue <- irc_messages.IrcSendMotdEnd()

	if ircClient.IsOsu {
		ircClient.packetQueue <- irc_messages.IrcSendPrivMsg("WaffleBot", "#osu", "Welcome to Waffle!")

		ircClient.ClientVersion = client_manager.ClientVersionOsuIrc
	} else {
		ircClient.ClientVersion = client_manager.ClientVersionIrc
	}

	client_manager.ClientManager.LockClientList()

	//Loop over every client which exists
	for _, currentClient := range client_manager.ClientManager.GetClientList() {
		//We already informed the new client, no need to do it again
		if currentClient.GetUserId() == int32(ircClient.UserData.UserID) {
			continue
		}

		relevantStats := ircClient.GetRelevantUserStats()

		//Inform client of our own existence
		currentClient.BanchoPresence(ircClient.UserData, relevantStats, 0)
		currentClient.BanchoOsuUpdate(relevantStats, ircClient.GetUserStatus())
	}

	client_manager.ClientManager.UnlockClientList()
	client_manager.ClientManager.RegisterClient(&ircClient)

	//Also register on the spectatable list
	spectator.ClientManager.RegisterClient(&ircClient)

	if ircClient.IsOsu {
		channelOsu, _ := chat.GetChannelByName("#osu")

		channelOsu.Join(&ircClient)

		ircClient.SendChannelNames(channelOsu)
		ircClient.joinedChannels["#osu"] = channelOsu
	}

	ircClient.lastPing = time.Now()
	ircClient.lastReceive = time.Now()
	ircClient.logonTime = time.Now()
	ircClient.continueRunning = true

	ctx, cancel := context.WithCancel(context.Background())

	ircClient.maintainCancel = cancel

	go ircClient.HandleIncoming()
	go ircClient.MaintainClient(ctx)
}

func (client *IrcClient) SendOffMessagesAndClose() {
	for len(client.packetQueue) != 0 {
		formatted, _ := (<-client.packetQueue).FormatMessage(client.Username)

		client.connection.Write([]byte(formatted))
	}

	client.connection.Close()
}

func (client *IrcClient) SendOffMessages() {
	for len(client.packetQueue) != 0 {
		formatted, _ := (<-client.packetQueue).FormatMessage(client.Username)

		client.connection.Write([]byte(formatted))
	}
}
