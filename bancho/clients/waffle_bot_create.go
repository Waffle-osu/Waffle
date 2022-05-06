package clients

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/packets"
	"Waffle/database"
	"Waffle/helpers"
	"sync"
	"time"
)

// CreateWaffleBot creates and brings WaffleBot to life
func CreateWaffleBot() {
	packetQueue := make(chan packets.BanchoPacket, 32)

	//Most of those are irrelevant cuz we aren't dealing with a real client
	clientInfo := ClientInformation{
		Timezone:       0,
		Version:        "Waffle",
		AllowCity:      false,
		OsuClientHash:  "https://github.com/Eeveelution/Waffle",
		MacAddressHash: "https://github.com/Eeveelution/Waffle",
	}

	fetchResult, user := database.UserFromDatabaseById(1)

	//If this happens, you either removed stuff from the DB or your MySQL stuff is wrong
	if fetchResult != 0 {
		helpers.Logger.Printf("[Bancho@WaffleBotCreate] ///////////// IMPORTANT //////////////")
		helpers.Logger.Printf("[Bancho@WaffleBotCreate] Failed to Find WaffleBot in Database!!")
		helpers.Logger.Printf("[Bancho@WaffleBotCreate] //////////////////////////////////////")

		return
	}

	statGetResult, osuStats := database.UserStatsFromDatabase(user.UserID, 0)
	statGetResult, taikoStats := database.UserStatsFromDatabase(user.UserID, 1)
	statGetResult, catchStats := database.UserStatsFromDatabase(user.UserID, 2)
	statGetResult, maniaStats := database.UserStatsFromDatabase(user.UserID, 3)

	//Makes the Rank not display in the client, good for distinguishing that this isn't a real player
	osuStats.Rank = 0
	taikoStats.Rank = 0
	catchStats.Rank = 0
	maniaStats.Rank = 0

	//If this happens, you either removed stuff from the DB or your MySQL stuff is wrong
	if statGetResult != 0 {
		helpers.Logger.Printf("[Bancho@WaffleBotCreate] //////////////// IMPORTANT /////////////////")
		helpers.Logger.Printf("[Bancho@WaffleBotCreate] Failed to Find WaffleBot stats in Database!!")
		helpers.Logger.Printf("[Bancho@WaffleBotCreate] ////////////////////////////////////////////")

		return
	}

	botClient := Client{
		//We don't need a connection because this is a local client
		connection:      nil,
		continueRunning: true,

		lastReceive: time.Now(),
		lastPing:    time.Now(),

		joinedChannels: make(map[string]*chat.Channel),
		awayMessage:    "",

		spectators:       make(map[int32]client_manager.OsuClient),
		spectatorMutex:   sync.Mutex{},
		spectatingClient: nil,

		isInLobby:         false,
		currentMultiLobby: nil,

		PacketQueue: packetQueue,

		UserData:   user,
		ClientData: clientInfo,
		Status: packets.StatusUpdate{
			Status:          packets.OsuStatusIdle,
			StatusText:      "Welcome to Waffle!",
			BeatmapChecksum: "No Map",
			CurrentMods:     0,
			Playmode:        packets.OsuGamemodeOsu,
			BeatmapId:       0,
		},
		OsuStats:    osuStats,
		TaikoStats:  taikoStats,
		CatchStats:  catchStats,
		ManiaStats:  maniaStats,
		FriendsList: []database.FriendEntry{},
	}

	client_manager.LockClientList()

	//Usually shouldn't matter because WaffleBot gets created the second bancho is and there's no way clients will connect this quick but ill keep it here
	for _, currentClient := range client_manager.GetClientList() {
		if currentClient.GetUserId() == int32(user.UserID) {
			continue
		}

		//Inform client of our own existence
		packets.BanchoSendUserPresence(currentClient.GetPacketQueue(), user, osuStats, clientInfo.Timezone)
		packets.BanchoSendOsuUpdate(currentClient.GetPacketQueue(), osuStats, botClient.Status)

		//Inform new client of the other client's existence
		packets.BanchoSendUserPresence(botClient.PacketQueue, currentClient.GetUserData(), currentClient.GetRelevantUserStats(), currentClient.GetClientTimezone())
		packets.BanchoSendOsuUpdate(botClient.PacketQueue, currentClient.GetRelevantUserStats(), currentClient.GetUserStatus())
	}

	client_manager.RegisterClient(&botClient)
	client_manager.UnlockClientList()

	//Since it has all permissions, it can join all channels it wants
	for _, channel := range chat.GetAvailableChannels() {
		channel.Join(&botClient)
	}

	//Starts Goroutines for handlig WaffleBot
	go botClient.WaffleBotMaintainClient()
	go botClient.WaffleBotHandleOutgoing()
}
