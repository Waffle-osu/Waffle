package waffle

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

type Client struct {
	connection net.Conn
}

func HandleNewClient(bancho *Bancho, connection net.Conn) {
	deadlineErr := connection.SetReadDeadline(time.Now().Add(5 * time.Second))

	if deadlineErr != nil {
		fmt.Printf("Failed to Configure 5 second read deadline.\n")
		return
	}

	textReader := bufio.NewReader(connection)

	username, err := textReader.ReadString('\n')
	password, err := textReader.ReadString('\n')
	userData, err := textReader.ReadString('\n')

	strings.Replace(username, "\r\n", "", -1)
	strings.Replace(password, "\r\n", "", -1)
	strings.Replace(userData, "\r\n", "", -1)

	if err != nil {
		fmt.Printf("Failed to read initial user data\n")
		return
	}

	fmt.Printf("Username: %s", username)
	fmt.Printf("Password: %s", password)
	fmt.Printf("UserData: %s", userData)

	connection.Close()
}
