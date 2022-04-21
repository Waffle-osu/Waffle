package clients

import (
	"sync"
)

var clients []*Client
var clientMutex sync.Mutex

func LockClientList() {
	clientMutex.Lock()
}

func UnlockClientList() {
	clientMutex.Unlock()
}

func GetClientList() []*Client {
	return clients
}

func GetClientByIndex(index int) *Client {
	return clients[index]
}

func GetAmountClients() int {
	return len(clients)
}

func RegisterClient(client *Client) {
	clients = append(clients, client)
}
