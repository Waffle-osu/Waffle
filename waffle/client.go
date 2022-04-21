package waffle

import (
	"Waffle/waffle/database"
	"Waffle/waffle/packets"
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ClientInformation struct {
	Timezone       int32
	Version        string
	AllowCity      bool
	OsuClientHash  string
	MacAddressHash string
}

type Client struct {
	Connection       net.Conn
	ClientData       ClientInformation
	BufReader        *bufio.Reader
	PacketQueue      *list.List
	PacketQueueMutex sync.Mutex
	UserData         database.User
	OsuStats         database.UserStats
	TaikoStats       database.UserStats
	CatchStats       database.UserStats
	ManiaStats       database.UserStats
}

func HandleNewClient(bancho *Bancho, connection net.Conn) {
	loginStartTime := time.Now()

	deadlineErr := connection.SetReadDeadline(time.Now().Add(5 * time.Second))

	if deadlineErr != nil {
		fmt.Printf("Failed to Configure 5 second read deadline.\n")
		return
	}

	textReader := bufio.NewReader(connection)

	username, readErr := textReader.ReadString('\n')
	password, readErr := textReader.ReadString('\n')
	userData, readErr := textReader.ReadString('\n')

	packetQueue := list.New()

	if readErr != nil {
		fmt.Printf("Failed to read initial user data\n")
		return
	}

	username = strings.Replace(username, "\r\n", "", -1)
	password = strings.Replace(password, "\r\n", "", -1)
	userData = strings.Replace(userData, "\r\n", "", -1)

	userDataSplit := strings.Split(userData, "|")

	if len(userDataSplit) != 4 {
		packets.BanchoSendLoginReply(packetQueue, packets.InvalidVersion)
		connection.Close()
		return
	}

	securityPartsSplit := strings.Split(userDataSplit[3], ":")

	timezone, convErr := strconv.Atoi(userDataSplit[1])

	if convErr != nil {
		packets.BanchoSendLoginReply(packetQueue, packets.InvalidVersion)
		connection.Close()
		return
	}

	clientInfo := ClientInformation{
		Version:        userDataSplit[0],
		Timezone:       int32(timezone),
		AllowCity:      userDataSplit[2] == "1",
		OsuClientHash:  securityPartsSplit[0],
		MacAddressHash: securityPartsSplit[1],
	}

	fetchResult, user := database.UserFromDatabaseByUsername(username)

	//No User Found
	if fetchResult == -1 {
		packets.BanchoSendLoginReply(packetQueue, packets.InvalidLogin)
		connection.Close()
		return
	} else if fetchResult == -2 {
		packets.BanchoSendLoginReply(packetQueue, packets.ServersideError)
		connection.Close()
		return
	}

	//Invalid Password
	if user.Password != password {
		packets.BanchoSendLoginReply(packetQueue, packets.InvalidLogin)
		connection.Close()
		return
	}

	//Banned
	if user.Banned == 1 {
		packets.BanchoSendLoginReply(packetQueue, packets.UserBanned)
		connection.Close()
		return
	}

	packets.BanchoSendLoginReply(packetQueue, int32(user.UserID))

	statGetResult, osuStats := database.UserStatsFromDatabase(user.UserID, 0)
	statGetResult, taikoStats := database.UserStatsFromDatabase(user.UserID, 0)
	statGetResult, catchStats := database.UserStatsFromDatabase(user.UserID, 0)
	statGetResult, maniaStats := database.UserStatsFromDatabase(user.UserID, 0)

	if statGetResult == -1 {
		//TODO: do a BanchoAnnounce to the user informing about the issue
		fmt.Printf("Uhh, user exists in users but not in stats")
		connection.Close()
		return
	} else if statGetResult == -2 {
		//TODO: do a BanchoAnnounce to the user informing about the issue
		connection.Close()
		return
	}

	osuClient := Client{
		Connection:  connection,
		PacketQueue: packetQueue,
		BufReader:   textReader,
		UserData:    user,
		ClientData:  clientInfo,
		OsuStats:    osuStats,
		TaikoStats:  taikoStats,
		CatchStats:  catchStats,
		ManiaStats:  maniaStats,
	}

	resetDeadlineErr := connection.SetReadDeadline(time.Time{})

	if resetDeadlineErr != nil {
		fmt.Printf("Failed to Configure 5 second read deadline.\n")
		return
	}

	packets.BanchoSendProtocolNegotiation(osuClient.PacketQueue)
	packets.BanchoSendLoginPermissions(osuClient.PacketQueue, user.Privileges)

	bancho.ClientMutex.Lock()
	bancho.Clients = append(bancho.Clients, &osuClient)
	bancho.ClientMutex.Unlock()

	fmt.Printf("Login for %s took %dus\n", username, time.Since(loginStartTime).Microseconds())
}

func (client *Client) HandleIncoming() {
	readBuffer := make([]byte, 4096)

	//Check if there's at least 1 packet header there
	availableBytes := client.BufReader.Buffered()

	if availableBytes > 0 {
		read, readErr := client.Connection.Read(readBuffer)

		if readErr != nil {
			return
		}

		//Get the bytes that were actually read
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

func (client *Client) SendOutgoing() {
	sendBuffer := new(bytes.Buffer)

	client.PacketQueueMutex.Lock()

	for retrievedPacket := client.PacketQueue.Front(); retrievedPacket != nil; retrievedPacket = retrievedPacket.Next() {
		packet := retrievedPacket.Value.(packets.BanchoPacket)

		fmt.Printf("Sending Packet %d\n", packet.PacketId)

		sendBuffer.Write(packet.GetBytes())

		client.PacketQueue.Remove(retrievedPacket)
	}

	client.PacketQueueMutex.Unlock()

	client.Connection.Write(sendBuffer.Bytes())
}
