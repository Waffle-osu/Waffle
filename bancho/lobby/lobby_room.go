package lobby

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/osu/base_packet_structures"
	"fmt"
	"sync"
	"time"
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

// JoinLobby is called when `client` joins the lobby
func JoinLobby(client LobbyClient) {
	LockClientList()

	//append to the client lists
	clientList = append(clientList, client)
	clientsById[client.GetUserId()] = client
	clientsByName[client.GetUserData().Username] = client

	//Inform everyone in the lobby that they joined
	for _, lobbyUser := range clientsById {
		client.BanchoLobbyJoin(lobbyUser.GetUserId())
		lobbyUser.BanchoLobbyJoin(client.GetUserId())
	}

	UnlockClientList()

	multiMutex.Lock()

	//Tell the new client of all the multiplayer matches that are going on
	for _, multiLobby := range multiLobbiesById {
		client.BanchoMatchNew(multiLobby.MatchInformation)
	}

	multiMutex.Unlock()
}

// PartLobby is called when `client` leaves the lobby
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
		lobbyUser.BanchoLobbyLeft(client.GetUserId())
	}

	UnlockClientList()
}

// BroadcastToLobby broadcasts a packet to everyone in the lobby
func BroadcastToLobby(packetFunction func(LobbyClient)) {
	LockClientList()

	for _, lobbyUser := range clientsById {
		packetFunction(lobbyUser)
	}

	UnlockClientList()
}

// CreateNewMultiMatch is responsible for creating a new Multiplayer Match
func CreateNewMultiMatch(match base_packet_structures.MultiplayerMatch, host LobbyClient, autojoinHost bool) *MultiplayerLobby {
	multiMutex.Lock()

	for i := 0; i != 65536; i++ {
		_, exists := multiLobbiesById[uint16(i)]

		if !exists {
			match.MatchId = uint16(i)
			break
		}
	}

	multiLobby := new(MultiplayerLobby)
	multiLobby.MatchId = fmt.Sprintf("%s-%s", host.GetUsername(), time.Now().Format("20060102150405"))

	//Set up the Chat channel #multiplayer
	multiLobby.MultiChannel = new(chat.Channel)
	multiLobby.MultiChannel.Name = "#multiplayer"
	multiLobby.MultiChannel.Description = ""
	multiLobby.MultiChannel.ReadPrivileges = chat.PrivilegesNormal
	multiLobby.MultiChannel.WritePrivileges = chat.PrivilegesNormal
	multiLobby.MultiChannel.Autojoin = false
	multiLobby.MultiChannel.Clients = []chat.ChatClient{}
	multiLobby.MultiChannel.ClientMutex = sync.Mutex{}

	//Set the match information, including host information
	multiLobby.MatchInformation = match
	multiLobby.MatchHost = host
	multiLobby.MatchInformation.HostId = host.GetUserId()

	//Make the host join the lobby, if specified
	if autojoinHost {
		host.JoinMatch(multiLobby, multiLobby.MatchInformation.GamePassword)
	}

	//Append lobby to the list
	multiLobbies = append(multiLobbies, multiLobby)
	multiLobbiesById[match.MatchId] = multiLobby

	//Tell everyone in the lobby that a new match has just been created
	BroadcastToLobby(func(client LobbyClient) {
		client.BanchoMatchNew(multiLobby.MatchInformation)
	})

	multiMutex.Unlock()

	return multiLobby
}

// RemoveMultiMatch gets called when a match gets disbanded
func RemoveMultiMatch(matchId uint16) {
	multiMutex.Lock()

	//Remove multi lobby from the multi list
	for index, value := range multiLobbies {
		if value.MatchInformation.MatchId == matchId {
			multiLobbies = append(multiLobbies[0:index], multiLobbies[index+1:]...)
		}
	}

	delete(multiLobbiesById, matchId)

	//Tell everyone in the lobby that the match no longer exists
	BroadcastToLobby(func(client LobbyClient) {
		client.BanchoMatchDisband(int32(matchId))
	})

	multiMutex.Unlock()
}

// GetMultiMatchById returns a multiplayer match given a match ID
func GetMultiMatchById(matchId uint16) *MultiplayerLobby {
	match, exists := multiLobbiesById[matchId]

	if !exists {
		return nil
	}

	return match
}

func GetMatchCount() int {
	return len(multiLobbies)
}
