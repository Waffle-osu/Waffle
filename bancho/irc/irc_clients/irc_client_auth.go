package irc_clients

import (
	"Waffle/bancho/irc/irc_messages"
	"Waffle/database"
	"Waffle/helpers"
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"net"
	"strings"
	"time"
)

func HandleNewIrcClient(connection net.Conn) {
	authBegin := time.Now()

	textReader := bufio.NewReader(connection)

	ircClient := IrcClient{
		connection: connection,
		reader:     textReader,
	}

	for ircClient.Username == "" && ircClient.Password == "" && time.Since(authBegin).Seconds() < 32 {
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
		//do things
		return
	}

	queryResult, foundUser := database.UserFromDatabaseById(uint64(userId))

	if queryResult != 0 {
		//do things

		return
	}

	ircClient.UserData = foundUser
}
