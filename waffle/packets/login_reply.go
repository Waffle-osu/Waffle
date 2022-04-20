package packets

import (
	"bytes"
	"encoding/binary"
	"net"
)

const (
	InvalidLogin          int32 = -1
	InvalidVersion        int32 = -2
	UserBanned            int32 = -3
	UnactivatedAccount    int32 = -4
	ServersideError       int32 = -5
	UnauthorizedTestBuild int32 = -6
)

func BanchoSendLoginReply(connection net.Conn, id int32) (bool, BanchoPacket) {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, id)

	loginReply := BanchoPacket{
		PacketId:   BanchoLoginReply,
		PacketSize: 4,
		PacketData: buf.Bytes(),
	}

	_, err := connection.Write(loginReply.GetBytes())

	if err != nil {
		return false, BanchoPacket{}
	}

	return false, BanchoPacket{}
}
