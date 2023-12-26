package spectator

import "Waffle/bancho/osu/base_packet_structures"

type SpectatorClient interface {
	// Retrieves this client's User ID
	GetUserId() int32
	// Retrieves the Username of the current client
	GetUsername() string

	// Sends the equivilant of a Spectator Join message.
	// Used to build a Spectator List
	BanchoSpectatorJoined(userId int32)
	// Sends the equivilant of a Spectator Leave message.
	// Used to build a Spectator List
	BanchoSpectatorLeft(userId int32)
	// Sends the equivilant of a Fellow Spectator Join message.
	// Used to build a Spectator List
	BanchoFellowSpectatorJoined(userId int32)
	// Sends the equivilant of a Fellow Spectator Leave message.
	// Used to build a Spectator List
	BanchoFellowSpectatorLeft(userId int32)
	// Sends the equivilant of a Spectator can't spectate message.
	// in osu! there's a seperate list for Spectators that don't have the map.
	BanchoSpectatorCantSpectate(userId int32)
	// Sends the equivilant of Spectator Replay Frames to the client.
	// This contains the next replay data of the client that this client is spectating
	BanchoSpectateFrames(frameBundle base_packet_structures.SpectatorFrameBundle)
}
