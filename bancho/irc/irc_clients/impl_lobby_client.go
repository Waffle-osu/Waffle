package irc_clients

import (
	"Waffle/bancho/chat"
	"Waffle/bancho/irc/irc_messages"
	"Waffle/bancho/lobby"
	"Waffle/bancho/osu/base_packet_structures"
)

func (ircClient *IrcClient) BanchoChannelRevoked(channel string) {
	ircClient.packetQueue <- irc_messages.IrcSendPart(ircClient.Username, channel)
}

func (ircClient *IrcClient) BanchoLobbyJoin(userId int32) {

}

func (ircClient *IrcClient) BanchoLobbyLeft(userId int32) {

}

func (ircClient *IrcClient) BanchoMatchNew(match base_packet_structures.MultiplayerMatch) {

}

func (ircClient *IrcClient) BanchoMatchUpdate(match base_packet_structures.MultiplayerMatch) {

}

func (ircClient *IrcClient) BanchoMatchStart(match base_packet_structures.MultiplayerMatch) {

}

func (ircClient *IrcClient) BanchoMatchDisband(matchId int32) {

}

func (ircClient *IrcClient) BanchoMatchTransferHost() {

}

func (ircClient *IrcClient) BanchoMatchAllPlayersLoaded() {

}

func (ircClient *IrcClient) BanchoMatchComplete() {

}

func (ircClient *IrcClient) BanchoMatchSkip() {

}

func (ircClient *IrcClient) BanchoMatchPlayerSkipped(slot int32) {

}

func (ircClient *IrcClient) BanchoMatchPlayerFailed(slot int32) {

}

func (ircClient *IrcClient) BanchoMatchScoreUpdate(scoreFrame base_packet_structures.ScoreFrame) {

}

func (ircClient *IrcClient) JoinMatch(match *lobby.MultiplayerLobby, password string) {

}

func (ircClient *IrcClient) LeaveCurrentMatch() {

}

func (ircClient *IrcClient) AssignMultiplayerLobby(lobby *lobby.MultiplayerLobby) {
	ircClient.currentMultiLobby = lobby
}

func (ircClient *IrcClient) GetMultiplayerLobby() *lobby.MultiplayerLobby {
	return ircClient.currentMultiLobby
}

func (ircClient *IrcClient) AddJoinedChannel(channel *chat.Channel) {
	ircClient.joinedChannels[channel.Name] = channel
}
