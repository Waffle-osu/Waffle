package clients

import (
	"Waffle/waffle/chat"
	"Waffle/waffle/client_manager"
	"Waffle/waffle/database"
	"Waffle/waffle/packets"
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

var guaranteedWorkingVersion = map[string]bool{
	"b1816.test":  true,
	"b1816.peppy": true,
	"b1816":       true,
	"b1815":       true,
}

func HandleNewClient(connection net.Conn) {
	loginStartTime := time.Now()

	deadlineErr := connection.SetReadDeadline(time.Now().Add(5 * time.Second))

	if deadlineErr != nil {
		fmt.Printf("Failed to Configure 5 second read deadline.\n")
		connection.Close()
		return
	}

	textReader := bufio.NewReader(connection)

	username, readErr := textReader.ReadString('\n')
	password, readErr := textReader.ReadString('\n')
	userData, readErr := textReader.ReadString('\n')

	packetQueue := make(chan packets.BanchoPacket, 32)

	if readErr != nil {
		fmt.Printf("Failed to read initial user data\n")
		connection.Close()
		return
	}

	username = strings.Replace(username, "\r\n", "", -1)
	password = strings.Replace(password, "\r\n", "", -1)
	userData = strings.Replace(userData, "\r\n", "", -1)

	userDataSplit := strings.Split(userData, "|")

	if len(userDataSplit) != 4 {
		packets.BanchoSendLoginReply(packetQueue, packets.InvalidVersion)
		go SendOffPacketsAndClose(connection, packetQueue)
		return
	}

	securityPartsSplit := strings.Split(userDataSplit[3], ":")

	timezone, convErr := strconv.Atoi(userDataSplit[1])

	if convErr != nil {
		packets.BanchoSendLoginReply(packetQueue, packets.InvalidVersion)
		go SendOffPacketsAndClose(connection, packetQueue)
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
		go SendOffPacketsAndClose(connection, packetQueue)
		return
	} else if fetchResult == -2 {
		packets.BanchoSendLoginReply(packetQueue, packets.ServersideError)
		go SendOffPacketsAndClose(connection, packetQueue)
		return
	}

	//Invalid Password
	if user.Password != password {
		packets.BanchoSendLoginReply(packetQueue, packets.InvalidLogin)
		go SendOffPacketsAndClose(connection, packetQueue)
		return
	}

	//Banned
	if user.Banned == 1 {
		packets.BanchoSendLoginReply(packetQueue, packets.UserBanned)
		go SendOffPacketsAndClose(connection, packetQueue)
		return
	}

	//Check for duplicate clients
	duplicateClient := client_manager.GetClientById(int32(user.UserID))

	if duplicateClient != nil {
		go func() {
			packets.BanchoSendAnnounce(duplicateClient.GetPacketQueue(), "Disconnecting because of another client conneting to your Account.")
			duplicateClient.CleanupClient()

			time.Sleep(2000)
			duplicateClient.Cut()
		}()
	}

	packets.BanchoSendLoginReply(packetQueue, int32(user.UserID))

	statGetResult, osuStats := database.UserStatsFromDatabase(user.UserID, 0)
	statGetResult, taikoStats := database.UserStatsFromDatabase(user.UserID, 1)
	statGetResult, catchStats := database.UserStatsFromDatabase(user.UserID, 2)
	statGetResult, maniaStats := database.UserStatsFromDatabase(user.UserID, 3)

	if statGetResult == -1 {
		packets.BanchoSendAnnounce(packetQueue, "A weird server-side fuckup occured, your stats don't exist yet your user does...")
		go SendOffPacketsAndClose(connection, packetQueue)
		return
	} else if statGetResult == -2 {
		packets.BanchoSendAnnounce(packetQueue, "A weird server-side fuckup occured, stats could not be loaded...")
		go SendOffPacketsAndClose(connection, packetQueue)
		return
	}

	client := Client{
		connection:      connection,
		lastPing:        time.Now(),
		lastReceive:     time.Now(),
		continueRunning: true,

		spectators: make(map[int32]client_manager.OsuClient),

		PacketQueue: packetQueue,

		UserData:   user,
		ClientData: clientInfo,
		OsuStats:   osuStats,
		TaikoStats: taikoStats,
		CatchStats: catchStats,
		ManiaStats: maniaStats,
		Status: packets.StatusUpdate{
			BeatmapChecksum: "",
			BeatmapId:       -1,
			CurrentMods:     0,
			Playmode:        packets.OsuGamemodeOsu,
			Status:          packets.OsuStatusIdle,
			StatusText:      user.Username + " has just logged in!",
		},
	}

	resetDeadlineErr := connection.SetReadDeadline(time.Time{})

	if resetDeadlineErr != nil {
		fmt.Printf("Failed to Configure 5 second read deadline.\n")
		go SendOffPacketsAndClose(connection, packetQueue)
	}

	packets.BanchoSendProtocolNegotiation(client.PacketQueue)
	packets.BanchoSendLoginPermissions(client.PacketQueue, user.Privileges|packets.UserPermissionsSupporter)
	packets.BanchoSendUserPresence(client.PacketQueue, user, osuStats, clientInfo.Timezone)
	packets.BanchoSendOsuUpdate(client.PacketQueue, osuStats, client.Status)

	client_manager.LockClientList()

	for i := 0; i != client_manager.GetAmountClients(); i++ {
		currentClient := client_manager.GetClientByIndex(i)

		if currentClient.GetUserId() == int32(user.UserID) {
			continue
		}

		//Inform client
		packets.BanchoSendUserPresence(currentClient.GetPacketQueue(), user, osuStats, clientInfo.Timezone)
		packets.BanchoSendOsuUpdate(currentClient.GetPacketQueue(), osuStats, client.Status)

		packets.BanchoSendUserPresence(client.PacketQueue, currentClient.GetUserData(), currentClient.GetRelevantUserStats(), currentClient.GetClientTimezone())
		packets.BanchoSendOsuUpdate(client.PacketQueue, currentClient.GetRelevantUserStats(), currentClient.GetUserStatus())
	}

	client_manager.RegisterClient(&client)
	client_manager.UnlockClientList()

	for _, channel := range chat.GetAvailableChannels() {
		if channel.Autojoin {
			if channel.Join(&client) {
				packets.BanchoSendChannelJoinSuccess(client.PacketQueue, channel.Name)
				client.joinedChannels = append(client.joinedChannels, channel)
			} else {
				packets.BanchoSendChannelRevoked(client.PacketQueue, channel.Name)
			}
		} else {
			packets.BanchoSendChannelAvailable(client.PacketQueue, channel.Name)
		}
	}

	working, recorded := guaranteedWorkingVersion[clientInfo.Version]

	if recorded == false {
		packets.BanchoSendAnnounce(client.PacketQueue, fmt.Sprintf("The osu! version %s has not yet been tested and may not work as intended! Unforseen problems may occur, report them to Furball if you can, depending on version it could be fixed.", clientInfo.Version))
	} else if working == false {
		packets.BanchoSendAnnounce(client.PacketQueue, fmt.Sprintf("The osu! version %s may not work as intended on waffle! Your experience may not be the best.", clientInfo.Version))
	} else {
		packets.BanchoSendAnnounce(client.PacketQueue, "Welcome to Waffle!")
	}

	fmt.Printf("%s successfully logged into Waffle using osu!%s\n", username, clientInfo.Version)
	fmt.Printf("Login took %.2fms\n", float64(time.Since(loginStartTime).Microseconds())/1000.0)

	go client.MaintainClient()
	go client.HandleIncoming()
	go client.SendOutgoing()
}

func SendOffPacketsAndClose(connection net.Conn, packetQueue chan packets.BanchoPacket) {
	for i := 0; i != len(packetQueue); i++ {
		connection.Write((<-packetQueue).GetBytes())
	}

	connection.Close()
}
