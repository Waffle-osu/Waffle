package lobby

import (
	"Waffle/waffle/chat"
	"Waffle/waffle/packets"
	"sync"
)

var clientList []LobbyClient
var clientsById map[int32]LobbyClient
var clientsByName map[string]LobbyClient
var clientMutex sync.Mutex

var multiLobbies []*MultiplayerLobby
var multiLobbiesById map[uint16]*MultiplayerLobby
var multiMutex sync.Mutex

func InitializeLobby() {
	clientsById = make(map[int32]LobbyClient)
	clientsByName = make(map[string]LobbyClient)
	multiLobbiesById = make(map[uint16]*MultiplayerLobby)
}

func LockClientList() {
	clientMutex.Lock()
}

func UnlockClientList() {
	clientMutex.Unlock()
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

	multiMutex.Lock()

	for _, multiLobby := range multiLobbiesById {
		packets.BanchoSendMatchNew(client.GetPacketQueue(), multiLobby.MatchInformation)
	}

	multiMutex.Unlock()
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

func BroadcastToLobby(packetFunction func(chan packets.BanchoPacket)) {
	LockClientList()

	for _, lobbyUser := range clientsById {
		packetFunction(lobbyUser.GetPacketQueue())
	}

	UnlockClientList()
}

func CreateNewMultiMatch(match packets.MultiplayerMatch, host LobbyClient) {
	multiMutex.Lock()

	for i := 0; i != 65536; i++ {
		_, exists := multiLobbiesById[uint16(i)]

		if exists == false {
			match.MatchId = uint16(i)
			break
		}
	}

	multiLobby := new(MultiplayerLobby)

	multiLobby.MultiChannel = new(chat.Channel)
	multiLobby.MultiChannel.Name = "#multiplayer"
	multiLobby.MultiChannel.Description = ""
	multiLobby.MultiChannel.AdminWrite = false
	multiLobby.MultiChannel.AdminRead = false
	multiLobby.MultiChannel.Autojoin = false
	multiLobby.MultiChannel.Clients = []chat.ChatClient{}
	multiLobby.MultiChannel.ClientMutex = sync.Mutex{}

	multiLobby.MatchInformation = match
	multiLobby.MatchHost = host
	multiLobby.MatchInformation.HostId = host.GetUserId()

	host.JoinMatch(multiLobby, multiLobby.MatchInformation.GamePassword)

	multiLobbies = append(multiLobbies, multiLobby)
	multiLobbiesById[match.MatchId] = multiLobby

	BroadcastToLobby(func(packetQueue chan packets.BanchoPacket) {
		packets.BanchoSendMatchNew(packetQueue, multiLobby.MatchInformation)
	})

	multiMutex.Unlock()
}

func RemoveMultiMatch(matchId uint16) {
	multiMutex.Lock()

	for index, value := range multiLobbies {
		if value.MatchInformation.MatchId == matchId {
			multiLobbies = append(multiLobbies[0:index], multiLobbies[index+1:]...)
		}
	}

	delete(multiLobbiesById, matchId)

	BroadcastToLobby(func(packetQueue chan packets.BanchoPacket) {
		packets.BanchoSendMatchDisband(packetQueue, int32(matchId))
	})

	multiMutex.Unlock()
}

func GetMultiMatchById(matchId uint16) *MultiplayerLobby {
	match, exists := multiLobbiesById[matchId]

	if exists == false {
		return nil
	}

	return match
}
