package client_manager

import (
	"sync"
)

var clientList []OsuClient
var clientsById map[int32]OsuClient
var clientsByName map[string]OsuClient
var clientMutex sync.Mutex

func InitializeClientManager() {
	clientsById = make(map[int32]OsuClient)
	clientsByName = make(map[string]OsuClient)
}

func LockClientList() {
	clientMutex.Lock()
}

func UnlockClientList() {
	clientMutex.Unlock()
}

func GetClientList() []OsuClient {
	return clientList
}

func GetClientByIndex(index int) OsuClient {
	return clientList[index]
}

func GetClientById(id int32) OsuClient {
	value, exists := clientsById[id]

	if exists == false {
		return nil
	}

	return value
}

func GetClientByName(username string) OsuClient {
	value, exists := clientsByName[username]

	if exists == false {
		return nil
	}

	return value
}

func GetAmountClients() int {
	return len(clientList)
}

func RegisterClient(client OsuClient) {
	clientList = append(clientList, client)
	clientsById[client.GetUserId()] = client
	clientsByName[client.GetUserData().Username] = client
}

func UnregisterClient(client OsuClient) {
	LockClientList()

	for index, value := range clientList {
		if value == client {
			clientList = append(clientList[0:index], clientList[index+1:]...)
		}
	}

	delete(clientsById, client.GetUserId())
	delete(clientsByName, client.GetUserData().Username)

	UnlockClientList()
}
