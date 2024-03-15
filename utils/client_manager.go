package utils

import "sync"

//The minimum we need from a client to enroll them in a Client list
type IdentifiableClient interface {
	// Retrieves this client's User ID
	GetUserId() int32
	// Retrieves the Username of the current client
	GetUsername() string
}

// Generic ClientManager for all purposes
type ClientManager[TClient IdentifiableClient] struct {
	clientList    []TClient
	clientsById   map[int32]TClient
	clientsByName map[string]TClient
	clientMutex   sync.Mutex
}

// Initializes the ClientManager
func (manager *ClientManager[TClient]) Initialize() {
	manager.clientList = []TClient{}
	manager.clientsById = make(map[int32]TClient)
	manager.clientsByName = make(map[string]TClient)
	manager.clientMutex = sync.Mutex{}
}

// Gets a client by their UserID
func (manager *ClientManager[TClient]) GetClientById(id int32) TClient {
	manager.clientMutex.Lock()
	defer manager.clientMutex.Unlock()

	return manager.clientsById[id]
}

// Gets a user by their Username
func (manager *ClientManager[TClient]) GetClientByName(username string) TClient {
	manager.clientMutex.Lock()
	defer manager.clientMutex.Unlock()

	return manager.clientsByName[username]
}

// Adds a client to the list
func (manager *ClientManager[TClient]) RegisterClient(client TClient) {
	manager.clientMutex.Lock()
	defer manager.clientMutex.Unlock()

	manager.clientList = append(manager.clientList, client)
	manager.clientsById[client.GetUserId()] = client
	manager.clientsByName[client.GetUsername()] = client
}

// Removes a client from the list
func (manager *ClientManager[TClient]) UnregisterClient(client TClient) {
	manager.clientMutex.Lock()
	defer manager.clientMutex.Unlock()

	for index, value := range manager.clientList {
		if value.GetUserId() == client.GetUserId() {
			manager.clientList = append(manager.clientList[0:index], manager.clientList[index+1:]...)
		}
	}

	delete(manager.clientsById, client.GetUserId())
	delete(manager.clientsByName, client.GetUsername())
}

// Gets the amount of clients in this list
func (manager *ClientManager[TClient]) GetClientCount() int {
	return len(manager.clientList)
}

// Locks the client list, useful for working with GetClientList as the list is guaranteed to not change.
// Don't forget to call UnlockClientList when you're done!
func (manager *ClientManager[TClient]) LockClientList() {
	manager.clientMutex.Lock()
}

// Unlocks the client list.
func (manager *ClientManager[TClient]) UnlockClientList() {
	manager.clientMutex.Unlock()
}

// Retreives the ClientList
func (manager *ClientManager[TClient]) GetClientList() []TClient {
	return manager.clientList
}
