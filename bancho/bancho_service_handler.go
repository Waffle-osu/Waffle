package bancho

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/client_manager"
	"Waffle/bancho/clients"
	"Waffle/bancho/lobby"
	"Waffle/bancho/misc"
	"Waffle/bancho/osu/b1815"
)

type BanchoService struct{}

func (service *BanchoService) Launch() {
	chat.InitializeChannels()                //Initializes Chat channels
	client_manager.InitializeClientManager() //Initializes the client manager
	lobby.InitializeLobby()                  //Initializes the multi lobby
	clients.WaffleBotInitializeCommands()    //Initializes Chat Commands
	misc.InitializeStatistics()              //Initializes Statistics
	b1815.InitializeCompatibilityLists()     //Initializes Client Compatibility lists

}

func (serice *BanchoService) Shutdown() {

}
