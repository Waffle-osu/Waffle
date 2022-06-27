package packets

import (
	"bytes"
	"encoding/binary"
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
	BanchoHandleIrcChangeUsername  uint16 = 9  //TODO: maybe do these once IRC's a thing?
	BanchoHandleIrcQuit            uint16 = 10 //TODO: maybe do these once IRC's a thing?
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
	BanchoGetAttention             uint16 = 23 //TODO: maybe once there's an admin panel or something? or maybe as a chat command for admins
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
	BanchoUnauthorized             uint16 = 62 //Unused
	OsuChannelJoin                 uint16 = 63
	BanchoChannelJoinSuccess       uint16 = 64
	BanchoChannelAvailable         uint16 = 65
	BanchoChannelRevoked           uint16 = 66
	BanchoChannelAvailableAutojoin uint16 = 67
	OsuBeatmapInfoRequest          uint16 = 68 //TODO: when you got all the maps in the DB you can do this
	BanchoBeatmapInfoReply         uint16 = 69 //TODO: when you got all the maps in the DB you can do this
	OsuMatchTransferHost           uint16 = 70
	BanchoLoginPermissions         uint16 = 71
	BanchoFriendsList              uint16 = 72
	OsuFriendsAdd                  uint16 = 73
	OsuFriendsRemove               uint16 = 74
	BanchoProtocolNegotiation      uint16 = 75
	BanchoTitleUpdate              uint16 = 76 //TODO: once site's a thing this could be used
	OsuMatchChangeTeam             uint16 = 77
	OsuChannelLeave                uint16 = 78
	OsuReceiveUpdates              uint16 = 79 //Unused
	BanchoMonitor                  uint16 = 80
	BanchoMatchPlayerSkipped       uint16 = 81
	OsuSetIrcAwayMessage           uint16 = 82
	BanchoUserPresence             uint16 = 83
	OsuUserStatsRequest            uint16 = 85
	BanchoRestart                  uint16 = 86

	BanchoHeaderSize int32 = 7
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

func ReadBanchoPacketHeader(packetBuffer *bytes.Buffer) (int, BanchoPacket, bool) {
	packet := BanchoPacket{}

	binary.Read(packetBuffer, binary.LittleEndian, &packet.PacketId)
	binary.Read(packetBuffer, binary.LittleEndian, &packet.PacketCompression)
	binary.Read(packetBuffer, binary.LittleEndian, &packet.PacketSize)

	if packet.PacketSize >= 32768 || packet.PacketId > 86 {
		return int(BanchoHeaderSize), packet, true
	}

	packet.PacketData = make([]byte, packet.PacketSize)

	binary.Read(packetBuffer, binary.LittleEndian, &packet.PacketData)

	return int(BanchoHeaderSize + packet.PacketSize), packet, false
}

func GetPacketName(packetId uint16) string {
	switch packetId {
	case OsuSendUserStatus:
		return "OsuSendUserStatus"
	case OsuSendIrcMessage:
		return "OsuSendIrcMessage"
	case OsuExit:
		return "OsuExit"
	case OsuRequestStatusUpdate:
		return "OsuRequestStatusUpdate"
	case OsuPong:
		return "OsuPong"
	case BanchoLoginReply:
		return "BanchoLoginReply"
	case BanchoSendMessage:
		return "BanchoSendMessage"
	case BanchoPing:
		return "BanchoPing"
	case BanchoHandleIrcChangeUsername:
		return "BanchoHandleIrcChangeUsername"
	case BanchoHandleIrcQuit:
		return "BanchoHandleIrcQuit"
	case BanchoHandleOsuUpdate:
		return "BanchoHandleOsuUpdate"
	case BanchoHandleOsuQuit:
		return "BanchoHandleOsuQuit"
	case BanchoSpectatorJoined:
		return "BanchoSpectatorJoined"
	case BanchoSpectatorLeft:
		return "BanchoSpectatorLeft"
	case BanchoSpectateFrames:
		return "BanchoSpectateFrames"
	case OsuStartSpectating:
		return "OsuStartSpectating"
	case OsuStopSpectating:
		return "OsuStopSpectating"
	case OsuSpectateFrames:
		return "OsuSpectateFrames"
	case OsuErrorReport:
		return "OsuErrorReport"
	case OsuCantSpectate:
		return "OsuCantSpectate"
	case BanchoSpectatorCantSpectate:
		return "BanchoSpectatorCantSpectate"
	case BanchoGetAttention:
		return "BanchoGetAttention"
	case BanchoAnnounce:
		return "BanchoAnnounce"
	case OsuSendIrcMessagePrivate:
		return "OsuSendIrcMessagePrivate"
	case BanchoMatchUpdate:
		return "BanchoMatchUpdate"
	case BanchoMatchNew:
		return "BanchoMatchNew"
	case BanchoMatchDisband:
		return "BanchoMatchDisband"
	case OsuLobbyPart:
		return "OsuLobbyPart"
	case OsuLobbyJoin:
		return "OsuLobbyJoin"
	case OsuMatchCreate:
		return "OsuMatchCreate"
	case OsuMatchJoin:
		return "OsuMatchJoin"
	case OsuMatchPart:
		return "OsuMatchPart"
	case BanchoLobbyJoin:
		return "BanchoLobbyJoin"
	case BanchoLobbyPart:
		return "BanchoLobbyPart"
	case BanchoMatchJoinSuccess:
		return "BanchoMatchJoinSuccess"
	case BanchoMatchJoinFail:
		return "BanchoMatchJoinFail"
	case OsuMatchChangeSlot:
		return "OsuMatchChangeSlot"
	case OsuMatchReady:
		return "OsuMatchReady"
	case OsuMatchLock:
		return "OsuMatchLock"
	case OsuMatchChangeSettings:
		return "OsuMatchChangeSettings"
	case BanchoFellowSpectatorJoined:
		return "BanchoFellowSpectatorJoined"
	case BanchoFellowSpectatorLeft:
		return "BanchoFellowSpectatorLeft"
	case OsuMatchStart:
		return "OsuMatchStart"
	case BanchoMatchStart:
		return "BanchoMatchStart"
	case OsuMatchScoreUpdate:
		return "OsuMatchScoreUpdate"
	case BanchoMatchScoreUpdate:
		return "BanchoMatchScoreUpdate"
	case OsuMatchComplete:
		return "OsuMatchComplete"
	case BanchoMatchTransferHost:
		return "BanchoMatchTransferHost"
	case OsuMatchChangeMods:
		return "OsuMatchChangeMods"
	case OsuMatchLoadComplete:
		return "OsuMatchLoadComplete"
	case BanchoMatchAllPlayersLoaded:
		return "BanchoMatchAllPlayersLoaded"
	case OsuMatchNoBeatmap:
		return "OsuMatchNoBeatmap"
	case OsuMatchNotReady:
		return "OsuMatchNotReady"
	case OsuMatchFailed:
		return "OsuMatchFailed"
	case BanchoMatchPlayerFailed:
		return "BanchoMatchPlayerFailed"
	case BanchoMatchComplete:
		return "BanchoMatchComplete"
	case OsuMatchHasBeatmap:
		return "OsuMatchHasBeatmap"
	case OsuMatchSkipRequest:
		return "OsuMatchSkipRequest"
	case BanchoMatchSkip:
		return "BanchoMatchSkip"
	case BanchoUnauthorized:
		return "BanchoUnauthorized"
	case OsuChannelJoin:
		return "OsuChannelJoin"
	case BanchoChannelJoinSuccess:
		return "BanchoChannelJoinSuccess"
	case BanchoChannelAvailable:
		return "BanchoChannelAvailable"
	case BanchoChannelRevoked:
		return "BanchoChannelRevoked"
	case BanchoChannelAvailableAutojoin:
		return "BanchoChannelAvailableAutojoin"
	case OsuBeatmapInfoRequest:
		return "OsuBeatmapInfoRequest"
	case BanchoBeatmapInfoReply:
		return "BanchoBeatmapInfoReply"
	case OsuMatchTransferHost:
		return "OsuMatchTransferHost"
	case BanchoLoginPermissions:
		return "BanchoLoginPermissions"
	case BanchoFriendsList:
		return "BanchoFriendsList"
	case OsuFriendsAdd:
		return "OsuFriendsAdd"
	case OsuFriendsRemove:
		return "OsuFriendsRemove"
	case BanchoProtocolNegotiation:
		return "BanchoProtocolNegotiation"
	case BanchoTitleUpdate:
		return "BanchoTitleUpdate"
	case OsuMatchChangeTeam:
		return "OsuMatchChangeTeam"
	case OsuChannelLeave:
		return "OsuChannelLeave"
	case OsuReceiveUpdates:
		return "OsuReceiveUpdates"
	case BanchoMonitor:
		return "BanchoMonitor"
	case BanchoMatchPlayerSkipped:
		return "BanchoMatchPlayerSkipped"
	case OsuSetIrcAwayMessage:
		return "OsuSetIrcAwayMessage"
	case BanchoUserPresence:
		return "BanchoUserPresence"
	case OsuUserStatsRequest:
		return "OsuUserStatsRequest"
	case BanchoRestart:
		return "BanchoRestart"
	}
	return ""
}
