package packets

import (
	"Waffle/database"
	"bytes"
	"encoding/binary"
)

func BanchoSendFriendsList(packetQueue chan BanchoPacket, friendsList []database.FriendEntry) {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, int16(len(friendsList)))

	for _, friend := range friendsList {
		binary.Write(buf, binary.LittleEndian, int32(friend.User2))
	}

	packetBytes := buf.Bytes()
	packetLength := len(packetBytes)

	packet := BanchoPacket{
		PacketId:          BanchoFriendsList,
		PacketCompression: 0,
		PacketSize:        int32(packetLength),
		PacketData:        packetBytes,
	}

	packetQueue <- packet
}
