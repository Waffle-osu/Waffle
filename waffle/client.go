package waffle

import (
	"Waffle/waffle/database"
	"Waffle/waffle/packets"
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strings"
	"time"
)

type Client struct {
	Connection  net.Conn
	BufReader   *bufio.Reader
	PacketQueue chan packets.BanchoPacket
	UserData    database.DatabaseUser
	OsuStats    database.DatabaseUserStats
	TaikoStats  database.DatabaseUserStats
	CatchStats  database.DatabaseUserStats
	ManiaStats  database.DatabaseUserStats
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
	if fetchResult < 0 {
		packets.BanchoSendLoginReply(connection, packets.InvalidLogin)
	}

	//Invalid Password
	if user.Password != password {
		packets.BanchoSendLoginReply(connection, packets.InvalidLogin)
	}

	//Banned
	if user.Banned == 1 {
		packets.BanchoSendLoginReply(connection, packets.UserBanned)
	}

	packets.BanchoSendLoginReply(connection, int32(user.UserID))

	fmt.Printf("Login for %s took %dus\n", username, time.Since(loginStartTime).Microseconds())

	osuClient := Client{
		Connection:  connection,
		PacketQueue: make(chan packets.BanchoPacket),
		BufReader:   textReader,
		UserData:    user,
	}

	bancho.ClientMutex.Lock()

	bancho.Clients = append(bancho.Clients, osuClient)

	bancho.ClientMutex.Unlock()

	resetDeadlineErr := connection.SetReadDeadline(time.Time{})

	if resetDeadlineErr != nil {
		fmt.Printf("Failed to Configure 5 second read deadline.\n")
		return
	}
}

func (client Client) HandleIncoming() {
	readBuffer := make([]byte, 4096)

	//Check if there's at least 1 packet header there
	_, peekErr := client.BufReader.Peek(packets.BanchoHeaderSize)

	if peekErr == nil {
		read, readErr := client.Connection.Read(readBuffer)

		if readErr != nil {
			return
		}

		packetBuffer := bytes.NewBuffer(readBuffer[:read])
		readIndex := 0

		for readIndex < read {
			read, packet := packets.ReadBanchoPacketHeader(packetBuffer)

			readIndex += read

			fmt.Printf("Read Packet ID: %d, of Size: %d, current readIndex: %d\n", packet.PacketId, packet.PacketSize, readIndex)

			//switch packet.PacketId {
			//
			//}
		}

	}
}
