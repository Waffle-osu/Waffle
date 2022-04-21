package clients

import (
	"sync"
)

var clientList []*Client
var clientsById map[int32]*Client
var clientMutex sync.Mutex

func InitializeClientManager() {
	clientsById = make(map[int32]*Client)
}

func LockClientList() {
	clientMutex.Lock()
}

func UnlockClientList() {
	clientMutex.Unlock()
}

func GetClientList() []*Client {
	return clientList
}

func GetClientByIndex(index int) *Client {
	return clientList[index]
}

func GetClientById(id int32) *Client {
	value, exists := clientsById[id]

	if exists == false {
		return nil
	}

	return value
}

func GetAmountClients() int {
	return len(clientList)
}

func RegisterClient(client *Client) {
	clientList = append(clientList, client)
	clientsById[int32(client.UserData.UserID)] = client
}

func UnregisterClient(client *Client) {
	LockClientList()

	for index, value := range clientList {
		if value == client {
			clientList = append(clientList[0:index], clientList[index+1:]...)
		}
	}

	delete(clientsById, int32(client.UserData.UserID))

	UnlockClientList()
}
