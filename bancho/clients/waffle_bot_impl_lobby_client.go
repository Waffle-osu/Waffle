package clients

import (
	"Waffle/bancho/lobby"
	"Waffle/bancho/osu/base_packet_structures"
)

func (client *WaffleBot) BanchoChannelRevoked(channel string) {

}

func (client *WaffleBot) BanchoLobbyJoin(userId int32) {

}

func (client *WaffleBot) BanchoLobbyLeft(userId int32) {

}

func (client *WaffleBot) BanchoMatchNew(match base_packet_structures.MultiplayerMatch) {

}

func (client *WaffleBot) BanchoMatchUpdate(match base_packet_structures.MultiplayerMatch) {

}

func (client *WaffleBot) BanchoMatchStart(match base_packet_structures.MultiplayerMatch) {

}

func (client *WaffleBot) BanchoMatchDisband(matchId int32) {

}

func (client *WaffleBot) BanchoMatchTransferHost() {

}

func (client *WaffleBot) BanchoMatchAllPlayersLoaded() {

}

func (client *WaffleBot) BanchoMatchComplete() {

}

func (client *WaffleBot) BanchoMatchSkip() {

}

func (client *WaffleBot) BanchoMatchPlayerSkipped(slot int32) {

}

func (client *WaffleBot) BanchoMatchPlayerFailed(slot int32) {

}

func (client *WaffleBot) BanchoMatchScoreUpdate(scoreFrame base_packet_structures.ScoreFrame) {

}

func (client *WaffleBot) JoinMatch(match *lobby.MultiplayerLobby, password string) {

}

func (client *WaffleBot) LeaveCurrentMatch() {

}

func (WaffleBot *WaffleBot) GetMultiplayerLobby() *lobby.MultiplayerLobby {
	return nil
}
