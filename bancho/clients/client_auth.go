package clients

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/packets"
	"Waffle/database"
	"Waffle/helpers"
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// All the currently tested versions and whether they work well or not
var guaranteedWorkingVersion = map[string]bool{
	//Official
	"b1816.test":  true,
	"b1816.peppy": true,
	"b1816":       true,
	"b1815":       true,
	"b1807":       true,
	"b1814":       true,
	"b1844.test":  true,

	//Unofficial
	"b1816modernized":                true,
	"b1816modernized.dev":            true,
	"b1816modernized.test":           true,
	"b1816modernized.ctbtest":        true,
	"b1816modernized-arcade":         true,
	"b1816modernized-arcade.dev":     true,
	"b1816modernized-arcade.test":    true,
	"b1816modernized-arcade.ctbtest": true,
}

// HandleNewClient handles a new connection
func HandleNewClient(connection net.Conn) {
	//Used to time how long a login takes
	loginStartTime := time.Now()

	deadlineErr := connection.SetReadDeadline(time.Now().Add(5 * time.Second))

	if deadlineErr != nil {
		helpers.Logger.Printf("[Bancho@Auth] Failed to Configure 5 second read deadline.\n")
		connection.Close()
		return
	}

	textReader := bufio.NewReader(connection)

	//Read everything the client gave us
	username, readErrUsername := textReader.ReadString('\n')
	password, readErrPassword := textReader.ReadString('\n')
	userData, readErrData := textReader.ReadString('\n')

	//Create a packet queue
	packetQueue := make(chan packets.BanchoPacket, 128)

	if readErrUsername != nil || readErrPassword != nil || readErrData != nil {
		helpers.Logger.Printf("[Bancho@Auth] Failed to read initial user data\n")
		connection.Close()
		return
	}

	//They have \r\n at the end, we trim that off
	username = strings.Replace(username, "\r", "", -1)
	password = strings.Replace(password, "\r", "", -1)
	userData = strings.Replace(userData, "\r", "", -1)

	username = strings.Replace(username, "\n", "", -1)
	password = strings.Replace(password, "\n", "", -1)
	userData = strings.Replace(userData, "\n", "", -1)

	//Start parsing the userData
	userDataSplit := strings.Split(userData, "|")

	//b1816 sends 4 components there, version|timezone|allow_city|security_parts
	if len(userDataSplit) != 4 {
		packets.BanchoSendLoginReply(packetQueue, packets.InvalidVersion)
		go SendOffPacketsAndClose(connection, packetQueue)
		return
	}

	//Parse security parts
	securityPartsSplit := strings.Split(userDataSplit[3], ":")

	//Parse timezone
	timezone, convErr := strconv.Atoi(userDataSplit[1])

	if convErr != nil {
		packets.BanchoSendLoginReply(packetQueue, packets.InvalidVersion)
		go SendOffPacketsAndClose(connection, packetQueue)
		return
	}

	//Construct Client information
	clientInfo := ClientInformation{
		Version:        userDataSplit[0],
		Timezone:       int32(timezone),
		AllowCity:      userDataSplit[2] == "1",
		OsuClientHash:  securityPartsSplit[0],
		MacAddressHash: securityPartsSplit[1],
	}

	//Try fetching the user from the Database
	fetchResult, user := database.UserFromDatabaseByUsername(username)

	//No User Found
	if fetchResult == -1 {
		packets.BanchoSendLoginReply(packetQueue, packets.InvalidLogin)
		go SendOffPacketsAndClose(connection, packetQueue)
		return
	} else if fetchResult == -2 {
		//Server failed to fetch the user
		packets.BanchoSendLoginReply(packetQueue, packets.ServersideError)
		go SendOffPacketsAndClose(connection, packetQueue)
		return
	}

	//Invalid Password
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
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
			duplicateClient.CleanupClient("Duplicate Client")

			//we wait for 2 seconds before cutting off the connection
			time.Sleep(2000)
			duplicateClient.Cut()
		}()
	}

	//Send successful login reply
	packets.BanchoSendLoginReply(packetQueue, int32(user.UserID))

	//Retrieve stats
	statGetResultOsu, osuStats := database.UserStatsFromDatabase(user.UserID, 0)
	statGetResultTaiko, taikoStats := database.UserStatsFromDatabase(user.UserID, 1)
	statGetResultCatch, catchStats := database.UserStatsFromDatabase(user.UserID, 2)
	statGetResultMania, maniaStats := database.UserStatsFromDatabase(user.UserID, 3)

	if statGetResultOsu == -1 || statGetResultTaiko == -1 || statGetResultCatch == -1 || statGetResultMania == -1 {
		packets.BanchoSendAnnounce(packetQueue, "A weird server-side fuckup occured, your stats don't exist yet your user does...")
		go SendOffPacketsAndClose(connection, packetQueue)
		return
	} else if statGetResultOsu == -2 || statGetResultTaiko == -2 || statGetResultCatch == -2 || statGetResultMania == -2 {
		packets.BanchoSendAnnounce(packetQueue, "A weird server-side fuckup occured, stats could not be loaded...")
		go SendOffPacketsAndClose(connection, packetQueue)
		return
	}

	//Retrieve friends list
	friendsResult, friendsList := database.FriendsGetFriendsList(user.UserID)

	if friendsResult != 0 {
		packets.BanchoSendAnnounce(packetQueue, "Friend List failed to load!")
	}

	//Send Friends list to client
	packets.BanchoSendFriendsList(packetQueue, friendsList)

	//Construct Client object
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
		FriendsList: friendsList,

		awayMessage:    "",
		joinedChannels: make(map[string]*chat.Channel),
	}

	resetDeadlineErr := connection.SetReadDeadline(time.Time{})

	if resetDeadlineErr != nil {
		helpers.Logger.Printf("[Bancho@Auth] Failed to Configure 5 second read deadline.\n")
		go SendOffPacketsAndClose(connection, packetQueue)
	}

	//Send Protocol negotiation aswell as information about itself
	packets.BanchoSendProtocolNegotiation(client.PacketQueue)
	packets.BanchoSendLoginPermissions(client.PacketQueue, user.Privileges|packets.UserPermissionsSupporter)
	packets.BanchoSendUserPresence(client.PacketQueue, user, osuStats, clientInfo.Timezone)
	packets.BanchoSendOsuUpdate(client.PacketQueue, osuStats, client.Status)

	client_manager.LockClientList()

	//Loop over every client which exists
	for _, currentClient := range client_manager.GetClientList() {
		//We already informed the new client, no need to do it again
		if currentClient.GetUserId() == int32(user.UserID) {
			continue
		}

		//Inform client of our own existence
		packets.BanchoSendUserPresence(currentClient.GetPacketQueue(), user, osuStats, clientInfo.Timezone)
		packets.BanchoSendOsuUpdate(currentClient.GetPacketQueue(), osuStats, client.Status)

		//Inform new client of the other client's existence
		packets.BanchoSendUserPresence(client.PacketQueue, currentClient.GetUserData(), currentClient.GetRelevantUserStats(), currentClient.GetClientTimezone())
		packets.BanchoSendOsuUpdate(client.PacketQueue, currentClient.GetRelevantUserStats(), currentClient.GetUserStatus())
	}

	//Register client in the client manager
	client_manager.RegisterClient(&client)
	client_manager.UnlockClientList()

	//Join all the channels we need/can
	for _, channel := range chat.GetAvailableChannels() {
		//No need to join admin channels, or stuff above our privilege
		if (channel.ReadPrivileges & user.Privileges) <= 0 {
			continue
		}

		if channel.Autojoin {
			if channel.Join(&client) {
				packets.BanchoSendChannelJoinSuccess(client.PacketQueue, channel.Name)
				client.joinedChannels[channel.Name] = channel
			} else {
				packets.BanchoSendChannelRevoked(client.PacketQueue, channel.Name)
			}
		} else {
			packets.BanchoSendChannelAvailable(client.PacketQueue, channel.Name)
		}
	}

	//Try getting info on the version the user's running
	working, recorded := guaranteedWorkingVersion[clientInfo.Version]

	if !recorded {
		packets.BanchoSendAnnounce(client.PacketQueue, fmt.Sprintf("The osu! version %s has not yet been tested and may not work as intended! Unforseen problems may occur, report them to Furball if you can, depending on version it could be fixed.", clientInfo.Version))
	} else if !working {
		packets.BanchoSendAnnounce(client.PacketQueue, fmt.Sprintf("The osu! version %s is tested and has been found to not work properly on Waffle! Your experience may not be the best.", clientInfo.Version))
	} else {
		packets.BanchoSendAnnounce(client.PacketQueue, "Welcome to Waffle!")
	}

	//Log some things
	helpers.Logger.Printf("[Bancho@Auth] %s successfully logged into Waffle using osu!%s\n", username, clientInfo.Version)
	helpers.Logger.Printf("[Bancho@Auth] Login took %.2fms\n", float64(time.Since(loginStartTime).Microseconds())/1000.0)

	//Start handlers
	go client.MaintainClient()
	go client.HandleIncoming()
	go client.SendOutgoing()
}

// SendOffPacketsAndClose sends off any remaining packets in the packet queue
func SendOffPacketsAndClose(connection net.Conn, packetQueue chan packets.BanchoPacket) {
	for len(packetQueue) != 0 {
		connection.Write((<-packetQueue).GetBytes())
	}

	connection.Close()
}
