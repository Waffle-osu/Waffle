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
