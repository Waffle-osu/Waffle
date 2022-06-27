package client_manager

import "Waffle/bancho/osu/b1815/packets"

//BroadcastPacket broadcasts a packet to everyone online
func BroadcastPacket(packetFunction func(packetQueue chan packets.BanchoPacket)) {
	for _, value := range clientList {
		packetFunction(value.GetPacketQueue())
	}
}
