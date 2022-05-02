package main

import (
	"Waffle/bancho"
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/clients"
	"Waffle/bancho/lobby"
	"Waffle/database"
	"Waffle/logger"
	"Waffle/web"
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

	_ = os.Mkdir(name, os.ModeDir)

	return true
}

func main() {
	EnsureDirectoryExists("logs")
	EnsureDirectoryExists("screenshots")
	EnsureDirectoryExists("release")

	logger.InitializeLogger()                //Initializes Logging, logs to both console and to a file
	chat.InitializeChannels()                //Initializes Chat channels
	client_manager.InitializeClientManager() //Initializes the client manager
	lobby.InitializeLobby()                  //Initializes the multi lobby
	clients.WaffleBotInitializeCommands()    //Initialize Chat Commands

	_, fileError := os.Stat(".env")

	if fileError != nil {
		database.Initialize("root", "root", "127.0.0.1:3306", "waffle")

		go func() {
			time.Sleep(time.Second * 2)

			logger.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")
			logger.Logger.Printf("[Initialization] //////////////////  First Run Advice  //////////////////\n")
			logger.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")
			logger.Logger.Printf("[Initialization] //   Run the osu!2011 Updater to configure waffle!!!  //\n")
			logger.Logger.Printf("[Initialization] //     You can set the MySQL Database and Location    //\n")
			logger.Logger.Printf("[Initialization] //      And more settings are likely coming soon!     //\n")
			logger.Logger.Printf("[Initialization] //                                                    //\n")
			logger.Logger.Printf("[Initialization] //      The Updater won't work properly until the     //\n")
			logger.Logger.Printf("[Initialization] //           Server is configured properly!           //\n")
			logger.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")
			logger.Logger.Printf("[Initialization] //            Or fill in the .env manually            //\n")
			logger.Logger.Printf("[Initialization] //        updater's cooler though  ¯\\_(ツ)_/¯         //\n")
			logger.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")
		}()
	} else {
		data, err := os.ReadFile(".env")

		if err != nil {
			logger.Logger.Fatalf("[Initialization] Failed to read configuration file, cannot start server!")
		}

		mySqlUsername := "root"
		mySqlPassword := "root"
		mySqlLocation := "127.0.0.1:3306"
		mySqlDatabase := "waffle"

		stringData := string(data)
		eachSetting := strings.Split(stringData, "\n")

		for _, iniEntry := range eachSetting {
			splitEntry := strings.Split(iniEntry, "=")

			if len(splitEntry) != 2 {
				continue
			}

			key := splitEntry[0]
			value := splitEntry[1]

			switch key {
			case "mysql_username":
				mySqlUsername = value
				break
			case "mysql_password":
				mySqlPassword = value
				break
			case "mysql_database":
				mySqlDatabase = value
				break
			case "mysql_location":
				mySqlLocation = value
				break
			}
		}

		database.Initialize(mySqlUsername, mySqlPassword, mySqlLocation, mySqlDatabase)
	}

	//Ensure all the updater items exist
	result, items := database.GetUpdaterItems()

	if result == -1 {
		logger.Logger.Printf("[Updater Checks] Failed to retrieve updater information!!!!!")
	}

	for _, item := range items {
		_, fileError := os.Stat("release/" + item.ServerFilename)

		if fileError != nil {
			logger.Logger.Printf("[Updater Checks] Updater Item File %s does not exist or cannot be accessed!\n", item.ServerFilename)
			logger.Logger.Printf("[Updater Checks] You can download the Updater Bundle here: https://eevee-sylveon.s-ul.eu/XqLHU708\n")
		}

		//Zip files will always have a mismatches hash, as they will be extracted client side
		if strings.HasSuffix(item.ServerFilename, ".zip") {
			continue
		}

		fileData, readErr := os.ReadFile("release/" + item.ServerFilename)

		if readErr != nil {
			logger.Logger.Printf("[Updater Checks] Updater Item File %s does not exist or cannot be accessed!\n", item.ServerFilename)
			logger.Logger.Printf("[Updater Checks] You can download the Updater Bundle here: https://eevee-sylveon.s-ul.eu/XqLHU708\n")
		}

		fileHash := md5.Sum(fileData)
		fileHashString := hex.EncodeToString(fileHash[:])

		if item.FileHash != fileHashString {
			logger.Logger.Printf("[Updater Checks] Updater Item File %s has mismatched MD5 Hashes!\n", item.ServerFilename)
			logger.Logger.Printf("[Updater Checks] Your hashes need to match in the database!\n")
			logger.Logger.Printf("[Updater Checks] You can download the Updater Bundle here: https://eevee-sylveon.s-ul.eu/XqLHU708\n")
		}
	}

	clients.CreateWaffleBot() //Creates WaffleBot

	go bancho.RunBancho()
	go web.RunOsuWeb()

	for {
		time.Sleep(2 * time.Second)
	}
}
