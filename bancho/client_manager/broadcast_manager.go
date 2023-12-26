package client_manager

// BroadcastPacket broadcasts a packet to everyone online
func BroadcastPacket(packetFunction func(client WaffleClient)) {
	for _, value := range clientList {
		packetFunction(value)
	}
}
