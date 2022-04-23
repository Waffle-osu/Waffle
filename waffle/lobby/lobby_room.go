package lobby

import (
	"Waffle/waffle/packets"
	"sync"
)

var clientList []LobbyClient
var clientsById map[int32]LobbyClient
var clientsByName map[string]LobbyClient
var clientMutex sync.Mutex

func InitializeLobby() {
	clientsById = make(map[int32]LobbyClient)
	clientsByName = make(map[string]LobbyClient)
}

func LockClientList() {
	clientMutex.Lock()
}

func UnlockClientList() {
	clientMutex.Unlock()
}

func GetClientList() []LobbyClient {
	return clientList
}

func GetClientByIndex(index int) LobbyClient {
	return clientList[index]
}

func GetClientById(id int32) LobbyClient {
	value, exists := clientsById[id]

	if exists == false {
		return nil
	}

	return value
}

func GetClientByName(username string) LobbyClient {
	value, exists := clientsByName[username]

	if exists == false {
		return nil
	}

	return value
}

func GetAmountClients() int {
	return len(clientList)
}

func JoinLobby(client LobbyClient) {
	LockClientList()

	clientList = append(clientList, client)
	clientsById[client.GetUserId()] = client
	clientsByName[client.GetUserData().Username] = client

	for _, lobbyUser := range clientsById {
		packets.BanchoSendLobbyJoin(client.GetPacketQueue(), lobbyUser.GetUserId())
		packets.BanchoSendLobbyJoin(lobbyUser.GetPacketQueue(), client.GetUserId())
	}

	UnlockClientList()
}

func PartLobby(client LobbyClient) {
	LockClientList()

	for index, value := range clientList {
		if value == client {
			clientList = append(clientList[0:index], clientList[index+1:]...)
		}
	}

	delete(clientsById, client.GetUserId())
	delete(clientsByName, client.GetUserData().Username)

	for _, lobbyUser := range clientsById {
		packets.BanchoSendLobbyPart(lobbyUser.GetPacketQueue(), client.GetUserId())
	}

	UnlockClientList()
}
