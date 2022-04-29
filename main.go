package main

import (
	"Waffle/bancho"
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/clients"
	"Waffle/bancho/lobby"
	"Waffle/database"
	"Waffle/web"
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

	clients.CreateWaffleBot() //Creates WaffleBot

	go bancho.RunBancho()
	go web.RunOsuWeb()

	for {
		time.Sleep(2 * time.Second)
	}
}
