package packets

import (
	"Waffle/waffle/database"
	"bytes"
	"encoding/binary"
)

func BanchoSendFriendsList(packetQueue chan BanchoPacket, friendsList []database.FriendEntry) {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, int32(len(friendsList)))

	for _, friend := range friendsList {
		binary.Write(buf, binary.LittleEndian, int32(friend.User2))
	}

	packetBytes := buf.Bytes()
	packetLength := len(packetBytes)

	packet := BanchoPacket{
		PacketId:          BanchoFriendsList, //TODO: change this out
		PacketCompression: 0,
		PacketSize:        int32(packetLength),
		PacketData:        packetBytes,
	}

	packetQueue <- packet
}
