package main

import (
	"Waffle/bancho"
	"Waffle/bancho/bot"
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/irc"
	"Waffle/bancho/lobby"
	"Waffle/bancho/misc"
	"Waffle/bancho/osu/b1815"
	"Waffle/bancho/spectator"
	"Waffle/config"
	"Waffle/database"
	"Waffle/helpers"
	"Waffle/scheduler"
	"crypto/md5"
	"encoding/hex"
	"os"
	"strings"
	"time"
)

func EnsureDirectoryExists(name string) bool {
	_, err := os.Stat(name)

	if err == nil {
		return false
	}

	_ = os.Mkdir(name, os.ModePerm)

	return true
}

func main() {
	EnsureDirectoryExists("logs")
	EnsureDirectoryExists("screenshots")
	EnsureDirectoryExists("release")
	EnsureDirectoryExists("replays")
	EnsureDirectoryExists("direct_thumbnails")
	EnsureDirectoryExists("mp3_previews")
	EnsureDirectoryExists("oszs")
	EnsureDirectoryExists("osus")
	EnsureDirectoryExists("avatars")
	EnsureDirectoryExists("bss_temp")

	helpers.InitializeLogger()               //Initializes Logging, logs to both console and to a file
	chat.InitializeChannels()                //Initializes Chat channels
	client_manager.InitializeClientManager() //Initializes the client manager
	spectator.InitializeClientManager()      //Initializes the spectator client manager
	lobby.InitializeLobby()                  //Initializes the multi lobby
	bot.WaffleBotInitializeCommands()        //Initializes Chat Commands
	misc.InitializeStatistics()              //Initializes Statistics
	b1815.InitializeCompatibilityLists()     //Initializes Client Compatibility lists
	config.ReadConfiguration()               //Initializes all Configurable things
	scheduler.InitializeJobScheduler()       //Initializes the Scheduler
	database.Initialize()                    //Initializes Database Connection and things
	database.InitializeMigrations()          //Initializes Database Migrations
	database.InitializeDatabaseVersion()     //Initializes the Current Database Version

	switch len(os.Args) {
	case 2:
		switch os.Args[1] {
		case "migrate":
			database.RunNecessaryMigrations()
			return
		case "more_data":
			GetAdditionalBeatmapInfo()
		}
	case 3:
		switch os.Args[1] {
		case "beatmap_versions":
			RunBeatmapClientVersionDetector(os.Args[2], false)
		case "beatmap_importer":
			BeatmapImporter(os.Args[2])
		case "osz_renamer":
			RenameOszs(os.Args[2])
		case "osu_mover":
			MoveOsuFiles(os.Args[2])
		case "migrate":
			database.RunNecessaryMigrations()
		}

		return
	case 4:
		switch os.Args[1] {
		case "beatmap_versions":
			RunBeatmapClientVersionDetector(os.Args[2], os.Args[3] == "write")
		}
	}

	//Ensure all the updater items exist
	result, items := database.UpdaterGetUpdaterItems()

	if result == -1 {
		helpers.Logger.Printf("[Updater Checks] Failed to retrieve updater information!!!!!")
	}

	for _, item := range items {
		_, fileError := os.Stat("release/" + item.ServerFilename)

		if fileError != nil {
			helpers.Logger.Printf("[Updater Checks] Updater Item File %s does not exist or cannot be accessed!\n", item.ServerFilename)
			helpers.Logger.Printf("[Updater Checks] You can download the Updater Bundle here: https://eevee-sylveon.s-ul.eu/XqLHU708\n")
		}

		//Zip files will always have a mismatches hash, as they will be extracted client side
		if strings.HasSuffix(item.ServerFilename, ".zip") {
			continue
		}

		fileData, readErr := os.ReadFile("release/" + item.ServerFilename)

		if readErr != nil {
			helpers.Logger.Printf("[Updater Checks] Updater Item File %s does not exist or cannot be accessed!\n", item.ServerFilename)
			helpers.Logger.Printf("[Updater Checks] You can download the Updater Bundle here: https://eevee-sylveon.s-ul.eu/XqLHU708\n")
		}

		fileHash := md5.Sum(fileData)
		fileHashString := hex.EncodeToString(fileHash[:])

		if item.FileHash != fileHashString {
			helpers.Logger.Printf("[Updater Checks] Updater Item File %s has mismatched MD5 Hashes!\n", item.ServerFilename)
			helpers.Logger.Printf("[Updater Checks] Your hashes need to match in the database!\n")
			helpers.Logger.Printf("[Updater Checks] You can download the Updater Bundle here: https://eevee-sylveon.s-ul.eu/XqLHU708\n")
		}
	}

	bot.CreateWaffleBot() //Creates WaffleBot

	go bancho.RunBancho()
	go RunWeb()
	go irc.RunIrcSSL()
	go irc.RunIrc()

	for {
		time.Sleep(2 * time.Second)
	}
}
