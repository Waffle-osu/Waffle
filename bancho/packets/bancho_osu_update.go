package packets

import (
	"Waffle/bancho/database"
	"bytes"
)

const (
	OsuStatusIdle         uint8 = 0
	OsuStatusAfk          uint8 = 1
	OsuStatusPlaying      uint8 = 2
	OsuStatusEditing      uint8 = 3
	OsuStatusModding      uint8 = 4
	OsuStatusMultiplayer  uint8 = 5
	OsuStatusWatching     uint8 = 6
	OsuStatusUnknown      uint8 = 7
	OsuStatusTesting      uint8 = 8
	OsuStatusSubmitting   uint8 = 9
	OsuStatusPaused       uint8 = 10
	OsuStatusLobby        uint8 = 11
	OsuStatusMultiplaying uint8 = 12
	OsuStatusOsuDirect    uint8 = 13
)

const (
	OsuGamemodeOsu   uint8 = 0
	OsuGamemodeTaiko uint8 = 1
	OsuGamemodeCatch uint8 = 2
	OsuGamemodeMania uint8 = 3
)

func BanchoSendOsuUpdate(packetQueue chan BanchoPacket, user database.UserStats, status StatusUpdate) {
	buf := new(bytes.Buffer)

	stats := OsuStats{
		UserId:      int32(user.UserID),
		Status:      status,
		RankedScore: int64(user.RankedScore),
		Accuracy:    user.Accuracy,
		Playcount:   int32(user.Playcount),
		TotalScore:  int64(user.TotalScore),
		Rank:        int32(user.Rank),
	}

	stats.WriteOsuStats(buf)

	packetBytes := buf.Bytes()
	packetLength := len(packetBytes)

	packet := BanchoPacket{
		PacketId:          BanchoHandleOsuUpdate,
		PacketCompression: 0,
		PacketSize:        int32(packetLength),
		PacketData:        packetBytes,
	}

	packetQueue <- packet
}
