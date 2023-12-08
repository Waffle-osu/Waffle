package b1815

import (
	"Waffle/helpers/serialization"
	"context"
)

type ClientEventListenType int32

type PacketEvent struct {
	PacketChannel chan serialization.BanchoPacket
	Context       context.Context
	Cancel        context.CancelFunc
}

// Adds event handler which fires every time a Packet is received
func (client *Client) OnPacket(handler func(packets chan serialization.BanchoPacket, ctx context.Context)) {
	channel := make(chan serialization.BanchoPacket, 128)
	ctx, cancel := context.WithCancel(context.Background())

	go handler(channel, ctx)

	client.packetListeners = append(client.packetListeners, PacketEvent{
		PacketChannel: channel,
		Context:       ctx,
		Cancel:        cancel,
	})
}
