package b1815

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/bancho/spectator"
	"Waffle/database"
	"Waffle/helpers"
	"Waffle/helpers/serialization"
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// All the currently tested versions and whether they work well or not
var guaranteedWorkingVersion = map[string]bool{
	/* Official Builds Starting Here */

	//b1816
	"b1816.test":  true,
	"b1816.peppy": true,
	"b1816":       true,
	//b1815
	"b1815":               true,
	"b1815.peppy":         true,
	"b1815.test":          true,
	"b1815.ctbtest":       true,
	"b1815arcade":         true,
	"b1815arcade.peppy":   true,
	"b1815arcade.test":    true,
	"b1815arcade.ctbtest": true,
	//b1807
	"b1807": true,
	//b1814
	"b1814": true,
	//b1844
	"b1844.test": true,
	//b20121119
	"b20121119":            false,
	"b20121119arcade":      false,
	"b20121119dev":         false,
	"b20121119public_test": false,

	/* Unofficial Builds Starting Here */

	//b1816 ported to FNA
	"b1816modernized":                true,
	"b1816modernized.dev":            true,
	"b1816modernized.test":           true,
	"b1816modernized.ctbtest":        true,
	"b1816modernized-arcade":         true,
	"b1816modernized-arcade.dev":     true,
	"b1816modernized-arcade.test":    true,
	"b1816modernized-arcade.ctbtest": true,
}

var knownIssuesList = map[string]string{}

func InitializeCompatibilityLists() {
	knownIssuesList["b20121119"] = "Multiplayer Matches cannot be created. Likely due to a difference in the MatchCreate packet, Leaderboards may fail to load, and the BanchoBeatmapInfoReply causes errors."
	knownIssuesList["b20121119arcade"] = knownIssuesList["b20121119"]
	knownIssuesList["b20121119dev"] = knownIssuesList["b20121119"]
	knownIssuesList["b20121119public_test"] = knownIssuesList["b20121119"]
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
	packetQueue := make(chan []byte, 128)

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
		packetQueue <- serialization.SendSerializableInt(serialization.BanchoLoginReply, serialization.InvalidVersion)

		go SendOffPacketsAndClose(connection, packetQueue)
		return
	}

	//Parse security parts
	securityPartsSplit := strings.Split(userDataSplit[3], ":")

	//Parse timezone
	timezone, convErr := strconv.Atoi(userDataSplit[1])

	if convErr != nil {
		packetQueue <- serialization.SendSerializableInt(serialization.BanchoLoginReply, serialization.InvalidVersion)

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
		packetQueue <- serialization.SendSerializableInt(serialization.BanchoLoginReply, serialization.InvalidLogin)
		go SendOffPacketsAndClose(connection, packetQueue)
		return
	} else if fetchResult == -2 {
		//Server failed to fetch the user
		packetQueue <- serialization.SendSerializableInt(serialization.BanchoLoginReply, serialization.ServersideError)
		go SendOffPacketsAndClose(connection, packetQueue)
		return
	}

	//Invalid Password
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		packetQueue <- serialization.SendSerializableInt(serialization.BanchoLoginReply, serialization.InvalidLogin)
		go SendOffPacketsAndClose(connection, packetQueue)
		return
	}

	//Banned
	if user.Banned == 1 {
		packetQueue <- serialization.SendSerializableInt(serialization.BanchoLoginReply, serialization.UserBanned)
		go SendOffPacketsAndClose(connection, packetQueue)
		return
	}

	//Check for duplicate clients
	duplicateClient := client_manager.ClientManager.GetClientById(int32(user.UserID))

	if duplicateClient != nil {
		go func() {
			packetQueue <- serialization.SendSerializableString(serialization.BanchoAnnounce, "Disconnecting because of another client conneting to your Account.")

			duplicateClient.CleanupClient("Duplicate Client")

			//we wait for 2 seconds before cutting off the connection
			time.Sleep(2000)
			duplicateClient.Cut()
		}()
	}

	//Send successful login reply
	packetQueue <- serialization.SendSerializableInt(serialization.BanchoLoginReply, int32(user.UserID))

	//Retrieve stats
	statGetResultOsu, osuStats := database.UserStatsFromDatabase(user.UserID, 0)
	statGetResultTaiko, taikoStats := database.UserStatsFromDatabase(user.UserID, 1)
	statGetResultCatch, catchStats := database.UserStatsFromDatabase(user.UserID, 2)
	statGetResultMania, maniaStats := database.UserStatsFromDatabase(user.UserID, 3)

	if statGetResultOsu == -1 || statGetResultTaiko == -1 || statGetResultCatch == -1 || statGetResultMania == -1 {
		packetQueue <- serialization.SendSerializableString(serialization.BanchoAnnounce, "A weird server-side fuckup occured, your stats don't exist yet your user does...")

		go SendOffPacketsAndClose(connection, packetQueue)
		return
	} else if statGetResultOsu == -2 || statGetResultTaiko == -2 || statGetResultCatch == -2 || statGetResultMania == -2 {
		packetQueue <- serialization.SendSerializableString(serialization.BanchoAnnounce, "A weird server-side fuckup occured, stats could not be loaded...")

		go SendOffPacketsAndClose(connection, packetQueue)
		return
	}

	//Retrieve friends list
	friendsResult, friendsList := database.FriendsGetFriendsList(user.UserID)

	if friendsResult != 0 {
		packetQueue <- serialization.SendSerializableString(serialization.BanchoAnnounce, "Friend List failed to load!")
	}

	//Send Friends list to client
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, int16(len(friendsList)))

	for _, friend := range friendsList {
		binary.Write(buf, binary.LittleEndian, int32(friend.User2))
	}

	packetQueue <- serialization.SendSerializableBytes(serialization.BanchoFriendsList, buf.Bytes())

	//Construct Client object
	client := Client{
		connection:      connection,
		lastReceive:     time.Now(),
		continueRunning: true,

		spectators: make(map[int32]spectator.SpectatorClient),

		PacketQueue: packetQueue,

		UserData:   user,
		ClientData: clientInfo,
		OsuStats:   osuStats,
		TaikoStats: taikoStats,
		CatchStats: catchStats,
		ManiaStats: maniaStats,
		Status: base_packet_structures.StatusUpdate{
			BeatmapChecksum: "",
			BeatmapId:       -1,
			CurrentMods:     0,
			Playmode:        serialization.OsuGamemodeOsu,
			Status:          serialization.OsuStatusIdle,
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
	presence := base_packet_structures.UserPresence{
		UserId:          int32(user.UserID),
		Username:        user.Username,
		AvatarExtension: 1,
		Timezone:        uint8(timezone),
		Country:         uint8(user.Country),
		City:            "",
		Permissions:     uint8(user.Privileges),
		Longitude:       0,
		Latitude:        0,
		Rank:            int32(osuStats.Rank),
	}

	stats := base_packet_structures.OsuStats{
		UserId:      int32(user.UserID),
		Status:      client.Status,
		RankedScore: int64(osuStats.RankedScore),
		Accuracy:    osuStats.Accuracy,
		Playcount:   int32(osuStats.Playcount),
		TotalScore:  int64(osuStats.TotalScore),
		Rank:        int32(osuStats.Rank),
	}

	packetQueue <- serialization.SendSerializableInt(serialization.BanchoProtocolNegotiation, 7)
	packetQueue <- serialization.SendSerializableInt(serialization.BanchoLoginPermissions, user.Privileges|serialization.UserPermissionsSupporter)
	packetQueue <- serialization.SendSerializable(serialization.BanchoUserPresence, presence)
	packetQueue <- serialization.SendSerializable(serialization.BanchoHandleOsuUpdate, stats)

	client_manager.ClientManager.LockClientList()

	//Loop over every client which exists
	for _, currentClient := range client_manager.ClientManager.GetClientList() {
		//We already informed the new client, no need to do it again
		if currentClient.GetUserId() == int32(user.UserID) {
			continue
		}

		//Inform client of our own existence
		currentClient.BanchoPresence(user, osuStats, clientInfo.Timezone)
		currentClient.BanchoOsuUpdate(osuStats, client.Status)

		//Inform new client of the other client's existence
		client.BanchoPresence(currentClient.GetUserData(), currentClient.GetRelevantUserStats(), currentClient.GetClientTimezone())
		client.BanchoOsuUpdate(currentClient.GetRelevantUserStats(), currentClient.GetUserStatus())
	}

	//Register client in the client manager
	client_manager.ClientManager.UnlockClientList()
	client_manager.ClientManager.RegisterClient(&client)

	//Also register on the spectatable list
	spectator.ClientManager.RegisterClient(&client)

	client.logonTime = time.Now()

	//Join all the channels we need/can
	for _, channel := range chat.GetAvailableChannels() {
		//No need to join admin channels, or stuff above our privilege
		if (channel.ReadPrivileges & user.Privileges) <= 0 {
			continue
		}

		if channel.Autojoin {
			if channel.Join(&client) {
				client.BanchoChannelJoinSuccess(channel.Name)

				client.joinedChannels[channel.Name] = channel
			} else {
				client.BanchoChannelRevoked(channel.Name)
			}
		} else {
			client.BanchoChannelAvailable(channel.Name)
		}
	}

	//Try getting info on the version the user's running
	working, recorded := guaranteedWorkingVersion[clientInfo.Version]

	if !recorded {
		client.BanchoAnnounce(fmt.Sprintf("The osu! version %s has not yet been tested and may not work as intended! Unforseen problems may occur, report them to Furball if you can, depending on version it could be fixed.", clientInfo.Version))
	} else if !working {
		client.BanchoAnnounce(fmt.Sprintf("The osu! version %s is tested and has been found to not work properly on Waffle! Your experience may not be the best.", clientInfo.Version))

		issues, hasKnownIssues := knownIssuesList[clientInfo.Version]

		if hasKnownIssues {
			client.BanchoAnnounce(fmt.Sprintf("The Client you're running on has these known issues:  %s", issues))
		}
	} else {
		client.BanchoAnnounce("Welcome to Waffle!")
	}

	//Log some things
	helpers.Logger.Printf("[Bancho@Auth] %s successfully logged into Waffle using osu!%s\n", username, clientInfo.Version)
	helpers.Logger.Printf("[Bancho@Auth] Login took %.2fms\n", float64(time.Since(loginStartTime).Microseconds())/1000.0)

	maintainCtx, cancelMaintain := context.WithCancel(context.Background())

	client.maintainCancel = cancelMaintain

	//Start handlers
	go client.MaintainClient(maintainCtx)
	go client.HandleIncoming()

	client.OnPacket(client.handlePackets) //Regular packet processing
	//client.OnPacket(client.waffleGuardPackets) //Waffle Guard processing for certaind detections
}

// SendOffPacketsAndClose sends off any remaining packets in the packet queue
func SendOffPacketsAndClose(connection net.Conn, packetQueue chan []byte) {
	for len(packetQueue) != 0 {
		connection.Write((<-packetQueue))
	}

	connection.Close()
}
