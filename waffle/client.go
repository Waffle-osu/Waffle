package waffle

import (
	"Waffle/waffle/database"
	"Waffle/waffle/packets"
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

type Client struct {
	Connection net.Conn
}

func HandleNewClient(bancho *Bancho, connection net.Conn) {
	loginStartTime := time.Now()

	deadlineErr := connection.SetReadDeadline(time.Now().Add(5 * time.Second))

	if deadlineErr != nil {
		fmt.Printf("Failed to Configure 5 second read deadline.\n")
		return
	}

	textReader := bufio.NewReader(connection)

	username, err := textReader.ReadString('\n')
	password, err := textReader.ReadString('\n')
	userData, err := textReader.ReadString('\n')

	username = strings.Replace(username, "\r\n", "", -1)
	password = strings.Replace(password, "\r\n", "", -1)
	userData = strings.Replace(userData, "\r\n", "", -1)

	if err != nil {
		fmt.Printf("Failed to read initial user data\n")
		return
	}

	fetchResult, user := database.UserFromDatabaseByUsername(username)

	//No User Found
	if fetchResult > 0 {
		packets.BanchoSendLoginReply(connection, -1)
	}

	//Invalid Password
	if user.Password != password {
		packets.BanchoSendLoginReply(connection, -1)
	}

	//Banned
	if user.Banned == 1 {
		packets.BanchoSendLoginReply(connection, -3)
	}

	packets.BanchoSendLoginReply(connection, int32(user.UserID))

	fmt.Printf("Login for %s took %dus\n", username, time.Since(loginStartTime).Microseconds())
}
