package client_manager

//BroadcastPacketOsu broadcasts a packet to everyone online
func BroadcastPacketOsu(packetFunction func(client WaffleClient)) {
	for _, value := range clientList {
		packetFunction(value)
	}
}
