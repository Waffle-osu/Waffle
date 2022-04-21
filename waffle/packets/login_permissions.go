package packets

import (
	"bytes"
	"container/list"
	"encoding/binary"
)

const (
	UserPermissionsRegular   = 1
	UserPermissionsBAT       = 2
	UserPermissionsSupporter = 4
	UserPermissionsFriend    = 8
)

func BanchoSendLoginPermissions(packetQueue *list.List, permissions int32) {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, permissions)

	packetBytes := buf.Bytes()
	packetLength := len(packetBytes)

	packet := BanchoPacket{
		PacketId:          BanchoLoginPermissions,
		PacketCompression: 0,
		PacketSize:        int32(packetLength),
		PacketData:        packetBytes,
	}

	packetQueue.PushBack(packet)
}
