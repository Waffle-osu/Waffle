package clients

import (
	"Waffle/waffle/chat"
	"Waffle/waffle/database"
	"Waffle/waffle/packets"
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func HandleNewClient(connection net.Conn) {
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

	packetQueue := make(chan packets.BanchoPacket, 32)

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
	statGetResult, taikoStats := database.UserStatsFromDatabase(user.UserID, 1)
	statGetResult, catchStats := database.UserStatsFromDatabase(user.UserID, 2)
	statGetResult, maniaStats := database.UserStatsFromDatabase(user.UserID, 3)

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

	client := Client{
		connection:      connection,
		lastPing:        time.Now(),
		lastReceive:     time.Now(),
		continueRunning: true,

		PacketQueue: packetQueue,

		UserData:   user,
		ClientData: clientInfo,
		OsuStats:   osuStats,
		TaikoStats: taikoStats,
		CatchStats: catchStats,
		ManiaStats: maniaStats,
		Status: packets.OsuStatus{
			BeatmapChecksum: "",
			BeatmapId:       -1,
			CurrentMods:     0,
			CurrentPlaymode: packets.OsuGamemodeOsu,
			CurrentStatus:   packets.OsuStatusIdle,
			StatusText:      user.Username + " has just logged in!",
		},
	}

	resetDeadlineErr := connection.SetReadDeadline(time.Time{})

	if resetDeadlineErr != nil {
		fmt.Printf("Failed to Configure 5 second read deadline.\n")
		return
	}

	packets.BanchoSendProtocolNegotiation(client.PacketQueue)
	packets.BanchoSendLoginPermissions(client.PacketQueue, user.Privileges)
	packets.BanchoSendUserPresence(client.PacketQueue, user, osuStats, clientInfo.Timezone)
	packets.BanchoSendOsuUpdate(client.PacketQueue, osuStats, client.Status)

	LockClientList()

	for i := 0; i != GetAmountClients(); i++ {
		currentClient := GetClientByIndex(i)

		if currentClient.UserData.UserID == user.UserID {
			continue
		}

		//Inform client
		packets.BanchoSendUserPresence(currentClient.PacketQueue, user, osuStats, clientInfo.Timezone)
		packets.BanchoSendOsuUpdate(currentClient.PacketQueue, osuStats, client.Status)

		var stats database.UserStats

		switch currentClient.Status.CurrentPlaymode {
		case packets.OsuGamemodeOsu:
			stats = currentClient.OsuStats
			break
		case packets.OsuGamemodeTaiko:
			stats = currentClient.TaikoStats
			break
		case packets.OsuGamemodeCatch:
			stats = currentClient.CatchStats
			break
		case packets.OsuGamemodeMania:
			stats = currentClient.ManiaStats
			break
		}

		packets.BanchoSendUserPresence(client.PacketQueue, currentClient.UserData, stats, currentClient.ClientData.Timezone)
		packets.BanchoSendOsuUpdate(client.PacketQueue, stats, currentClient.Status)
	}

	RegisterClient(&client)
	UnlockClientList()

	if chat.TryJoinChannel(&client, "#osu") {
		packets.BanchoSendChannelJoinSuccess(client.PacketQueue, "#osu")

		client.joinedChannels = append(client.joinedChannels, "#osu")
	}

	if chat.TryJoinChannel(&client, "#announce") {
		packets.BanchoSendChannelJoinSuccess(client.PacketQueue, "#announce")

		client.joinedChannels = append(client.joinedChannels, "#announce")
	}

	fmt.Printf("Login for %s took %dus\n", username, time.Since(loginStartTime).Microseconds())

	go client.MaintainClient()
	go client.HandleIncoming()
	go client.SendOutgoing()
}
