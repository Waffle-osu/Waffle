package main

import (
	"Waffle/bancho"
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/clients"
	"Waffle/bancho/lobby"
	"Waffle/database"
	"Waffle/web"
	"crypto/md5"
	"encoding/hex"
	"fmt"
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
	EnsureDirectoryExists("screenshots")
	EnsureDirectoryExists("release")

	chat.InitializeChannels()                //Initializes Chat channels
	client_manager.InitializeClientManager() //Initializes the client manager
	lobby.InitializeLobby()                  //Initializes the multi lobby
	clients.WaffleBotInitializeCommands()    //Initialize Chat Commands

	_, fileError := os.Stat(".env")

	if fileError != nil {
		database.Initialize("root", "root", "127.0.0.1:3306", "waffle")

		go func() {
			time.Sleep(time.Second * 2)

			fmt.Printf("////////////////////////////////////////////////////////\n")
			fmt.Printf("//////////////////  First Run Advice  //////////////////\n")
			fmt.Printf("////////////////////////////////////////////////////////\n")
			fmt.Printf("//   Run the osu!2011 Updater to configure waffle!!!  //\n")
			fmt.Printf("//     You can set the MySQL Database and Location    //\n")
			fmt.Printf("//      And more settings are likely coming soon!     //\n")
			fmt.Printf("//                                                    //\n")
			fmt.Printf("//      The Updater won't work properly until the     //\n")
			fmt.Printf("//           Server is configured properly!           //\n")
			fmt.Printf("////////////////////////////////////////////////////////\n")
			fmt.Printf("//            Or fill in the .env manually            //\n")
			fmt.Printf("//        updater's cooler though  ¯\\_(ツ)_/¯         //\n")
			fmt.Printf("////////////////////////////////////////////////////////\n")
		}()
	} else {
		data, err := os.ReadFile(".env")

		if err != nil {
			fmt.Printf("Failed to read configuration file, cannot start server!")
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
		fmt.Printf("Failed to retrieve updater information!!!!!")
	}

	for _, item := range items {
		_, fileError := os.Stat("release/" + item.ServerFilename)

		if fileError != nil {
			fmt.Printf("Updater Item File %s does not exist or cannot be accessed!\n", item.ServerFilename)
			fmt.Printf("You can download the Updater Bundle here: https://eevee-sylveon.s-ul.eu/XqLHU708\n")
		}

		if strings.HasSuffix(item.ServerFilename, ".zip") {
			continue
		}
		fileData, readErr := os.ReadFile("release/" + item.ServerFilename)

		if readErr != nil {
			fmt.Printf("Updater Item File %s does not exist or cannot be accessed!\n", item.ServerFilename)
			fmt.Printf("You can download the Updater Bundle here: https://eevee-sylveon.s-ul.eu/XqLHU708\n")
		}

		fileHash := md5.Sum(fileData)
		fileHashString := hex.EncodeToString(fileHash[:])

		if item.FileHash != fileHashString {
			fmt.Printf("Updater Item File %s has mismatched MD5 Hashes!\n", item.ServerFilename)
			fmt.Printf("Your hashes need to match in the database!\n")
			fmt.Printf("You can download the Updater Bundle here: https://eevee-sylveon.s-ul.eu/XqLHU708\n")
		}
	}

	clients.CreateWaffleBot() //Creates WaffleBot

	go bancho.RunBancho()
	go web.RunOsuWeb()

	for {
		time.Sleep(2 * time.Second)
	}
}
