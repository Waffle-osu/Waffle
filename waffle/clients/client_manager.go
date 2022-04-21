package clients

import (
	"sync"
)

var clientList []*Client
var clientMutex sync.Mutex

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

func GetAmountClients() int {
	return len(clientList)
}

func RegisterClient(client *Client) {
	clientList = append(clientList, client)
}

func UnregisterClient(client *Client) {
	LockClientList()

	for index, value := range clientList {
		if value == client {
			clientList = append(clientList[0:index], clientList[index+1:]...)
		}
	}

	UnlockClientList()
}
