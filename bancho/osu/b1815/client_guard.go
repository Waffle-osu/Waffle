package b1815

import (
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/helpers"
	"Waffle/helpers/serialization"
	"bytes"
	"context"
	"time"
)

type WaffleGuardContext struct {
	lastStatusUpdate base_packet_structures.StatusUpdate
	playingStart     time.Time
}

func (client *Client) waffleGuardPackets(packetChannel chan serialization.BanchoPacket, ctx context.Context) {
	var guardCtx *WaffleGuardContext = &client.waffleGuardContext

	for {
		select {
		case packet := <-packetChannel:
			//Packet data reader, only contains the packet data
			packetDataReader := bytes.NewBuffer(packet.PacketData)

			switch packet.PacketId {
			case serialization.OsuSendUserStatus:
				statusUpdate := base_packet_structures.ReadStatusUpdate(packetDataReader)

				wasPlaying := guardCtx.lastStatusUpdate.Status == serialization.OsuStatusPlaying || guardCtx.lastStatusUpdate.Status == serialization.OsuStatusMultiplaying
				statusChanged := guardCtx.lastStatusUpdate.Status != statusUpdate.Status

				if wasPlaying && statusChanged {
					playStop := time.Now()
					timePlayed := playStop.Sub(guardCtx.playingStart)

					helpers.Guard.Printf("%s played %s for %.2f seconds", client.UserData.Username, guardCtx.lastStatusUpdate.StatusText, timePlayed.Seconds())
				}

				if statusUpdate.Status == serialization.OsuStatusPlaying || statusUpdate.Status == serialization.OsuStatusMultiplaying {
					helpers.Guard.Printf("%s is now playing %s; tracking...", client.UserData.Username, statusUpdate.StatusText)

					guardCtx.lastStatusUpdate = statusUpdate
					guardCtx.playingStart = time.Now()
				}
			}
		case <-ctx.Done():

		default:
		}
	}
}
