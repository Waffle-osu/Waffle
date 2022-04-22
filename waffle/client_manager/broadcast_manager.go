package client_manager

import "Waffle/waffle/packets"

func BroadcastPacket(packetFunction func(packetQueue chan packets.BanchoPacket)) {
	for _, value := range clientList {
		packetFunction(value.GetPacketQueue())
	}
}
