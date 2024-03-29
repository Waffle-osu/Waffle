package b1815

import (
	"Waffle/bancho/misc"
	"Waffle/helpers"
	"Waffle/helpers/packets"
	"bytes"
	"context"
	"time"
)

// HandleIncoming handles things coming from the osu! client
func (client *Client) HandleIncoming() {
	//make a 32kb Buffer to read stuff
	readBuffer := make([]byte, 32768)

	for client.continueRunning {
		read, readErr := client.connection.Read(readBuffer)

		if readErr != nil {
			//We don't clean up as we may not need to
			continue
		}

		go func() {
			misc.StatsRecvLock.Lock()
			misc.StatsBytesRecieved += uint64(read)
			misc.StatsRecvLock.Unlock()
		}()

		//Update last receive time, this is used to check for timeouts
		client.lastReceive = time.Now()

		//Get the bytes that were actually read
		packetBuffer := bytes.NewBuffer(readBuffer[:read])
		//Index into the buffer, so we read every packet that we have
		readIndex := 0

		for readIndex < read {
			read, packet, failedRead := packets.ReadBanchoPacketHeader(packetBuffer)

			readIndex += read

			if failedRead {
				continue
			}

			//Unused packet
			if packet.PacketId == 79 {
				continue
			}

			for _, packetReceiver := range client.packetListeners {
				packetReceiver.PacketChannel <- packet
			}
		}
	}
}

// MaintainClient is looping every second, sending out pings and handles timeouts
func (client *Client) MaintainClient(ctx context.Context) {
	pingTicker := time.NewTicker(PingTimeout * time.Second)
	receiveTicker := time.NewTicker(ReceiveTimeout * time.Second)

	for {
		select {
		//Maintains the client until the client is cancelled.
		case <-ctx.Done():
			//We close in MaintainClient instead of in CleanupClient to avoid possible double closes, causing panics
			helpers.Logger.Printf("[Bancho@Handling] Closed %s's Packet Queue", client.UserData.Username)
			close(client.PacketQueue)

			pingTicker.Stop()
			receiveTicker.Stop()
			return
		//Sends off outgoing packets
		case packet := <-client.PacketQueue:
			sendBytes := len(packet)

			go func() {
				misc.StatsSendLock.Lock()
				misc.StatsBytesSent += uint64(sendBytes)
				misc.StatsSendLock.Unlock()
			}()

			_, err := client.connection.Write(packet)

			if err != nil {
				helpers.Logger.Printf(err.Error())
			}
		// Sends pings
		case <-pingTicker.C:
			client.BanchoPing()
		//Closes the client on timeout
		case <-receiveTicker.C:
			lastReceiveUnix := client.lastReceive.Unix()
			unixNow := time.Now().Unix()

			if lastReceiveUnix+ReceiveTimeout <= unixNow {
				client.CleanupClient("Client Timed out.")

				client.continueRunning = false
			}
		}
	}
}
