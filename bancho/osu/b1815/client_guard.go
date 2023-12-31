package b1815

import (
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/helpers/packets"
	"context"
	"time"
)

type WaffleGuardContext struct {
	lastStatusUpdate base_packet_structures.StatusUpdate
	playingStart     time.Time
}

func (client *Client) waffleGuardPackets(packetChannel chan packets.BanchoPacket, ctx context.Context) {
	// var guardCtx *WaffleGuardContext = &client.waffleGuardContext

	for {
		select {
		// case packet := <-packetChannel:
		//Packet data reader, only contains the packet data
		// packetDataReader := bytes.NewBuffer(packet.PacketData)

		// switch packet.PacketId {
		// }
		case <-ctx.Done():
			return
		}
	}
}
