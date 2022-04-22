package packets

import (
	"bytes"
	"encoding/binary"
	"io"
)

const (
	OsuSendUserStatus              uint16 = 0
	OsuSendIrcMessage              uint16 = 1
	OsuExit                        uint16 = 2
	OsuRequestStatusUpdate         uint16 = 3
	OsuPong                        uint16 = 4
	BanchoLoginReply               uint16 = 5
	BanchoSendMessage              uint16 = 7
	BanchoPing                     uint16 = 8
	BanchoHandleIrcChangeUsername  uint16 = 9
	BanchoHandleIrcQuit            uint16 = 10
	BanchoHandleOsuUpdate          uint16 = 11
	BanchoHandleOsuQuit            uint16 = 12
	BanchoSpectatorJoined          uint16 = 13
	BanchoSpectatorLeft            uint16 = 14
	BanchoSpectateFrames           uint16 = 15
	OsuStartSpectating             uint16 = 16
	OsuStopSpectating              uint16 = 17
	OsuSpectateFrames              uint16 = 18
	OsuErrorReport                 uint16 = 20
	OsuCantSpectate                uint16 = 21
	BanchoSpectatorCantSpectate    uint16 = 22
	BanchoGetAttention             uint16 = 23
	BanchoAnnounce                 uint16 = 24
	OsuSendIrcMessagePrivate       uint16 = 25
	BanchoMatchUpdate              uint16 = 26
	BanchoMatchNew                 uint16 = 27
	BanchoMatchDisband             uint16 = 28
	OsuLobbyPart                   uint16 = 29
	OsuLobbyJoin                   uint16 = 30
	OsuMatchCreate                 uint16 = 31
	OsuMatchJoin                   uint16 = 32
	OsuMatchPart                   uint16 = 33
	BanchoLobbyJoin                uint16 = 34
	BanchoLobbyPart                uint16 = 35
	BanchoMatchJoinSuccess         uint16 = 36
	BanchoMatchJoinFail            uint16 = 37
	OsuMatchChangeSlot             uint16 = 38
	OsuMatchReady                  uint16 = 39
	OsuMatchLock                   uint16 = 40
	OsuMatchChangeSettings         uint16 = 41
	BanchoFellowSpectatorJoined    uint16 = 42
	BanchoFellowSpectatorLeft      uint16 = 43
	OsuMatchStart                  uint16 = 44
	BanchoMatchStart               uint16 = 46
	OsuMatchScoreUpdate            uint16 = 47
	BanchoMatchScoreUpdate         uint16 = 48
	OsuMatchComplete               uint16 = 49
	BanchoMatchTransferHost        uint16 = 50
	OsuMatchChangeMods             uint16 = 51
	OsuMatchLoadComplete           uint16 = 52
	BanchoMatchAllPlayersLoaded    uint16 = 53
	OsuMatchNoBeatmap              uint16 = 54
	OsuMatchNotReady               uint16 = 55
	OsuMatchFailed                 uint16 = 56
	BanchoMatchPlayerFailed        uint16 = 57
	BanchoMatchComplete            uint16 = 58
	OsuMatchHasBeatmap             uint16 = 59
	OsuMatchSkipRequest            uint16 = 60
	BanchoMatchSkip                uint16 = 61
	BanchoUnauthorized             uint16 = 62
	OsuChannelJoin                 uint16 = 63
	BanchoChannelJoinSuccess       uint16 = 64
	BanchoChannelAvailable         uint16 = 65
	BanchoChannelRevoked           uint16 = 66
	BanchoChannelAvailableAutojoin uint16 = 67
	OsuBeatmapInfoRequest          uint16 = 68
	BanchoBeatmapInfoReply         uint16 = 69
	OsuMatchTransferHost           uint16 = 70
	BanchoLoginPermissions         uint16 = 71
	BanchoFriendsList              uint16 = 72
	OsuFriendsAdd                  uint16 = 73
	OsuFriendsRemove               uint16 = 74
	BanchoProtocolNegotiation      uint16 = 75
	BanchoTitleUpdate              uint16 = 76
	OsuMatchChangeTeam             uint16 = 77
	OsuChannelLeave                uint16 = 78
	OsuReceiveUpdates              uint16 = 79
	BanchoMonitor                  uint16 = 80
	BanchoMatchPlayerSkipped       uint16 = 81
	OsuSetIrcAwayMessage           uint16 = 82
	BanchoUserPresence             uint16 = 83
	OsuUserStatsRequest            uint16 = 85
	BanchoRestart                  uint16 = 86

	BanchoHeaderSize int = 7
)

type BanchoPacket struct {
	PacketId          uint16
	PacketCompression int8
	PacketSize        int32
	PacketData        []byte
}

func (packet BanchoPacket) GetBytes() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, packet.PacketId)
	binary.Write(buf, binary.LittleEndian, packet.PacketCompression)
	binary.Write(buf, binary.LittleEndian, packet.PacketSize)

	if packet.PacketSize != 0 {
		binary.Write(buf, binary.LittleEndian, packet.PacketData)
	}

	return buf.Bytes()
}

func ReadBanchoPacketHeader(packetBuffer *bytes.Buffer) (int, BanchoPacket) {
	packet := BanchoPacket{}

	binary.Read(packetBuffer, binary.LittleEndian, &packet.PacketId)
	binary.Read(packetBuffer, binary.LittleEndian, &packet.PacketCompression)
	binary.Read(packetBuffer, binary.LittleEndian, &packet.PacketSize)

	packet.PacketData = make([]byte, packet.PacketSize)

	binary.Read(packetBuffer, binary.LittleEndian, &packet.PacketData)

	return int(7 + packet.PacketSize), packet
}

func WriteBanchoString(value string) []byte {
	if value == "" {
		return []byte{0}
	}

	var length int
	var i int = len(value)
	var ulebBytes []byte

	if i == 0 {
		ulebBytes = []byte{0}
	}

	for i > 0 {
		ulebBytes = append(ulebBytes, 0)
		ulebBytes[length] = byte(i & 0x7F)
		i >>= 7
		if i != 0 {
			ulebBytes[length] |= 0x80
		}
		length++
	}

	returnBytes := []byte{11}
	returnBytes = append(returnBytes, ulebBytes...)
	returnBytes = append(returnBytes, []byte(value)...)

	return returnBytes
}

func ReadBanchoString(reader io.Reader) []byte {
	bytes := make([]byte, 1)

	reader.Read(bytes)

	if bytes[0] != 11 {
		return []byte{}
	}

	var shift uint
	var lastByte byte
	var total int

	for {
		b := make([]byte, 1)
		reader.Read(b)
		lastByte = b[0]
		total |= (int(lastByte&0x7F) << shift)
		if lastByte&0x80 == 0 {
			break
		}
		shift += 7
	}

	bytes = make([]byte, total)

	reader.Read(bytes)

	return bytes
}
