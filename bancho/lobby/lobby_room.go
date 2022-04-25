package lobby

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/packets"
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

// JoinLobby is called when `client` joins the lobby
func JoinLobby(client LobbyClient) {
	LockClientList()

	//append to the client lists
	clientList = append(clientList, client)
	clientsById[client.GetUserId()] = client
	clientsByName[client.GetUserData().Username] = client

	//Inform everyone in the lobby that they joined
	for _, lobbyUser := range clientsById {
		packets.BanchoSendLobbyJoin(client.GetPacketQueue(), lobbyUser.GetUserId())
		packets.BanchoSendLobbyJoin(lobbyUser.GetPacketQueue(), client.GetUserId())
	}

	UnlockClientList()

	multiMutex.Lock()

	//Tell the new client of all the multiplayer matches that are going on
	for _, multiLobby := range multiLobbiesById {
		packets.BanchoSendMatchNew(client.GetPacketQueue(), multiLobby.MatchInformation)
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
		packets.BanchoSendLobbyPart(lobbyUser.GetPacketQueue(), client.GetUserId())
	}

	UnlockClientList()

	//TODO@(Furball): i don't know whether this bug exists but it might very well exist,
	//              : what happens when a user leaves the lobby and a match disbands,
	//              : does the client still remember the lobby or does it disappear?
}

// BroadcastToLobby broadcasts a packet to everyone in the lobby
func BroadcastToLobby(packetFunction func(chan packets.BanchoPacket)) {
	LockClientList()

	for _, lobbyUser := range clientsById {
		packetFunction(lobbyUser.GetPacketQueue())
	}

	UnlockClientList()
}

// CreateNewMultiMatch is responsible for creating a new Multiplayer Match
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

	//Make the host join the lobby
	host.JoinMatch(multiLobby, multiLobby.MatchInformation.GamePassword)

	//Append lobby to the list
	multiLobbies = append(multiLobbies, multiLobby)
	multiLobbiesById[match.MatchId] = multiLobby

	//Tell everyone in the lobby that a new match has just been created
	BroadcastToLobby(func(packetQueue chan packets.BanchoPacket) {
		packets.BanchoSendMatchNew(packetQueue, multiLobby.MatchInformation)
	})

	multiMutex.Unlock()
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
	BroadcastToLobby(func(packetQueue chan packets.BanchoPacket) {
		packets.BanchoSendMatchDisband(packetQueue, int32(matchId))
	})

	multiMutex.Unlock()
}

// GetMultiMatchById returns a multiplayer match given a match ID
func GetMultiMatchById(matchId uint16) *MultiplayerLobby {
	match, exists := multiLobbiesById[matchId]

	if exists == false {
		return nil
	}

	return match
}
