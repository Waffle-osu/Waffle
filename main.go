package main

import (
	"Waffle/bancho"
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/clients"
	"Waffle/bancho/lobby"
	"Waffle/bancho/misc"
	"Waffle/database"
	"Waffle/helpers"
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
	EnsureDirectoryExists("replays")
	EnsureDirectoryExists("direct_thumbnails")
	EnsureDirectoryExists("mp3_previews")
	EnsureDirectoryExists("oszs")

	helpers.InitializeLogger()               //Initializes Logging, logs to both console and to a file
	chat.InitializeChannels()                //Initializes Chat channels
	client_manager.InitializeClientManager() //Initializes the client manager
	lobby.InitializeLobby()                  //Initializes the multi lobby
	clients.WaffleBotInitializeCommands()    //Initializes Chat Commands
	misc.InitializeStatistics()              //Initializes Statistics

	_, fileError := os.Stat(".env")

	if fileError != nil {
		database.Initialize("root", "root", "127.0.0.1:3306", "waffle")

		go func() {
			time.Sleep(time.Second * 2)

			defaultConfig := "mysql_username=root\nmysql_password=root\nmysql_location=127.0.0.1:3306\nmysql_database=waffle"

			writeErr := os.WriteFile(".env", []byte(defaultConfig), 0644)

			if writeErr != nil {
				helpers.Logger.Printf("[Initialization] Failed to create default configuration!\n\n")
				return
			}

			helpers.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")
			helpers.Logger.Printf("[Initialization] //////////////////  First Run Advice  //////////////////\n")
			helpers.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")
			helpers.Logger.Printf("[Initialization] // No .env file detected in the Waffle.exe directory! //\n")
			helpers.Logger.Printf("[Initialization] //    This file stores important configuration for    //\n")
			helpers.Logger.Printf("[Initialization] //      The server, such as Database Credentials,     //\n")
			helpers.Logger.Printf("[Initialization] //                                                    //\n")
			helpers.Logger.Printf("[Initialization] // A .env file with default settings has been created //\n")
			helpers.Logger.Printf("[Initialization] //      Please change the settings as necessary       //\n")
			helpers.Logger.Printf("[Initialization] //                                                    //\n")
			helpers.Logger.Printf("[Initialization] //    Explanation to all the keys in the .env file:   //\n")
			helpers.Logger.Printf("[Initialization] //                                                    //\n")
			helpers.Logger.Printf("[Initialization] // mysql_location: Location of the MySQL Server       //\n")
			helpers.Logger.Printf("[Initialization] // mysql_database: Name of the Database to use        //\n")
			helpers.Logger.Printf("[Initialization] // mysql_username: Under which user to log in         //\n")
			helpers.Logger.Printf("[Initialization] // mysql_password: Password for said user             //\n")
			helpers.Logger.Printf("[Initialization] ////////////////////////////////////////////////////////\n")
		}()
	} else {
		data, err := os.ReadFile(".env")

		if err != nil {
			helpers.Logger.Fatalf("[Initialization] Failed to read configuration file, cannot start server!")
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
			case "mysql_password":
				mySqlPassword = value
			case "mysql_database":
				mySqlDatabase = value
			case "mysql_location":
				mySqlLocation = value
			}
		}

		database.Initialize(mySqlUsername, mySqlPassword, mySqlLocation, mySqlDatabase)
	}

	if len(os.Args) == 3 {
		switch os.Args[1] {
		case "beatmap_importer":
			BeatmapImporter(os.Args[2])
		case "osz_renamer":
			RenameOszs(os.Args[2])
		}

		return
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

	clients.CreateWaffleBot() //Creates WaffleBot

	go bancho.RunBancho()
	go web.RunOsuWeb()

	for {
		time.Sleep(2 * time.Second)
	}
}
