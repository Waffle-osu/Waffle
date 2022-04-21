package packets

import (
	"bytes"
	"encoding/binary"
)

const (
	InvalidLogin          int32 = -1
	InvalidVersion        int32 = -2
	UserBanned            int32 = -3
	UnactivatedAccount    int32 = -4
	ServersideError       int32 = -5
	UnauthorizedTestBuild int32 = -6
)

func BanchoSendLoginReply(packetQueue chan BanchoPacket, id int32) {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, id)

	loginReply := BanchoPacket{
		PacketId:   BanchoLoginReply,
		PacketSize: 4,
		PacketData: buf.Bytes(),
	}

	packetQueue <- loginReply
}
