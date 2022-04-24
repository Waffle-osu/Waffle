package clients

import (
	"Waffle/waffle/chat"
	"Waffle/waffle/client_manager"
	"Waffle/waffle/database"
	"Waffle/waffle/packets"
	"fmt"
	"sync"
	"time"
)

func CreateWaffleBot() {
	packetQueue := make(chan packets.BanchoPacket, 32)

	clientInfo := ClientInformation{
		Timezone:       0,
		Version:        "Waffle",
		AllowCity:      false,
		OsuClientHash:  "https://github.com/Eeveelution/Waffle",
		MacAddressHash: "https://github.com/Eeveelution/Waffle",
	}

	fetchResult, user := database.UserFromDatabaseById(1)

	if fetchResult != 0 {
		fmt.Printf("///////////// IMPORTANT //////////////")
		fmt.Printf("Failed to Find WaffleBot in Database!!")
		fmt.Printf("//////////////////////////////////////")

		return
	}

	statGetResult, osuStats := database.UserStatsFromDatabase(user.UserID, 0)
	statGetResult, taikoStats := database.UserStatsFromDatabase(user.UserID, 1)
	statGetResult, catchStats := database.UserStatsFromDatabase(user.UserID, 2)
	statGetResult, maniaStats := database.UserStatsFromDatabase(user.UserID, 3)

	osuStats.Rank = 0
	taikoStats.Rank = 0
	catchStats.Rank = 0
	maniaStats.Rank = 0

	if statGetResult != 0 {
		fmt.Printf("//////////////// IMPORTANT /////////////////")
		fmt.Printf("Failed to Find WaffleBot stats in Database!!")
		fmt.Printf("////////////////////////////////////////////")

		return
	}

	botClient := Client{
		connection:      nil,
		continueRunning: true,

		lastReceive: time.Now(),
		lastPing:    time.Now(),

		joinedChannels: []*chat.Channel{},
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

	for i := 0; i != client_manager.GetAmountClients(); i++ {
		currentClient := client_manager.GetClientByIndex(i)

		if currentClient.GetUserId() == int32(user.UserID) {
			continue
		}

		//Inform client
		packets.BanchoSendUserPresence(currentClient.GetPacketQueue(), user, osuStats, clientInfo.Timezone)
		packets.BanchoSendOsuUpdate(currentClient.GetPacketQueue(), osuStats, botClient.Status)

		packets.BanchoSendUserPresence(botClient.PacketQueue, currentClient.GetUserData(), currentClient.GetRelevantUserStats(), currentClient.GetClientTimezone())
		packets.BanchoSendOsuUpdate(botClient.PacketQueue, currentClient.GetRelevantUserStats(), currentClient.GetUserStatus())
	}

	client_manager.RegisterClient(&botClient)
	client_manager.UnlockClientList()

	for _, channel := range chat.GetAvailableChannels() {
		channel.Join(&botClient)
	}

	go botClient.WaffleBotMaintainClient()
	go botClient.WaffleBotHandleOutgoing()
}
