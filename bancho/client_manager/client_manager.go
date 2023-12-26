package client_manager

import (
	"Waffle/common"
)

/*
	For the record, client_manager.ClientManager is fucking stupid
	But, because this is golang, there's nothing I can do about it!
	I'd put it in bancho, but bancho needs b1815, and if b1815 needs bancho
	boom circular dependency. Until I figure out a better solution
	this comedy stays
*/

var ClientManager common.ClientManager[WaffleClient] = common.ClientManager[WaffleClient]{}

// InitializeClientManager initializes the ClientManager
func InitializeClientManager() {
	ClientManager.Initialize()
}

// BroadcastPacket broadcasts a packet to everyone online
func BroadcastPacket(packetFunction func(client WaffleClient)) {
	ClientManager.LockClientList()

	for _, value := range ClientManager.GetClientList() {
		packetFunction(value)
	}

	ClientManager.UnlockClientList()
}
