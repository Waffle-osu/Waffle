package bot

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/database"
	"Waffle/helpers"
	"Waffle/helpers/packets"
	"sync"
	"time"
)

var WaffleBotInstance *WaffleBot

// CreateWaffleBot creates and brings WaffleBot to life
func CreateWaffleBot() {
	fetchResult, user := database.UserFromDatabaseById(1)

	//If this happens, you either removed stuff from the DB or your MySQL stuff is wrong
	if fetchResult != 0 {
		helpers.Logger.Printf("[Bancho@WaffleBotCreate] ///////////// IMPORTANT //////////////")
		helpers.Logger.Printf("[Bancho@WaffleBotCreate] Failed to Find WaffleBot in Database!!")
		helpers.Logger.Printf("[Bancho@WaffleBotCreate] //////////////////////////////////////")

		return
	}

	statGetResultOsu, osuStats := database.UserStatsGetWaffleBot(0)
	statGetResultTaiko, taikoStats := database.UserStatsGetWaffleBot(1)
	statGetResultCatch, catchStats := database.UserStatsGetWaffleBot(2)
	statGetResultMania, maniaStats := database.UserStatsGetWaffleBot(3)

	//Makes the Rank not display in the client, good for distinguishing that this isn't a real player
	osuStats.Rank = 0
	taikoStats.Rank = 0
	catchStats.Rank = 0
	maniaStats.Rank = 0

	//If this happens, you either removed stuff from the DB or your MySQL stuff is wrong
	if statGetResultOsu != 0 || statGetResultTaiko != 0 || statGetResultCatch != 0 || statGetResultMania != 0 {
		helpers.Logger.Printf("[Bancho@WaffleBotCreate] //////////////// IMPORTANT /////////////////")
		helpers.Logger.Printf("[Bancho@WaffleBotCreate] Failed to Find WaffleBot stats in Database!!")
		helpers.Logger.Printf("[Bancho@WaffleBotCreate] Please create a user called WaffleBot under the User ID of 1")
		helpers.Logger.Printf("[Bancho@WaffleBotCreate] ////////////////////////////////////////////")

		return
	}

	botClient := WaffleBot{
		continueRunning: true,

		lastReceive: time.Now(),
		lastPing:    time.Now(),
		logonTime:   time.Now(),

		joinedChannels: make(map[string]*chat.Channel),
		awayMessage:    "",

		spectators:       make(map[int32]client_manager.WaffleClient),
		spectatorMutex:   sync.Mutex{},
		spectatingClient: nil,

		isInLobby:         false,
		currentMultiLobby: nil,

		UserData: user,
		Status: base_packet_structures.StatusUpdate{
			Status:          packets.OsuStatusIdle,
			StatusText:      "Welcome to Waffle!",
			BeatmapChecksum: "No Map",
			CurrentMods:     0,
			Playmode:        packets.OsuGamemodeOsu,
			BeatmapId:       0,
		},
		OsuStats:   osuStats,
		TaikoStats: taikoStats,
		CatchStats: catchStats,
		ManiaStats: maniaStats,
	}

	WaffleBotInstance = &botClient

	client_manager.ClientManager.LockClientList()

	//Usually shouldn't matter because WaffleBot gets created the second bancho is and there's no way clients will connect this quick but ill keep it here
	for _, currentClient := range client_manager.ClientManager.GetClientList() {
		if currentClient.GetUserId() == int32(user.UserID) {
			continue
		}

		//Inform client of our own existence
		currentClient.BanchoPresence(user, osuStats, 0)
		currentClient.BanchoOsuUpdate(osuStats, botClient.Status)

		//Inform new client of the other client's existence
		botClient.BanchoPresence(currentClient.GetUserData(), currentClient.GetRelevantUserStats(), currentClient.GetClientTimezone())
		botClient.BanchoOsuUpdate(currentClient.GetRelevantUserStats(), currentClient.GetUserStatus())
	}

	client_manager.ClientManager.UnlockClientList()
	client_manager.ClientManager.RegisterClient(&botClient)

	//Since it has all permissions, it can join all channels it wants
	for _, channel := range chat.GetAvailableChannels() {
		channel.Join(&botClient)
	}

	//Starts Goroutines for handling WaffleBot
	go botClient.WaffleBotMaintainClient()
}
