package packets

import (
	"Waffle/waffle/database"
	"bytes"
	"encoding/binary"
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

type OsuStatus struct {
	BeatmapChecksum string
	BeatmapId       int32
	CurrentMods     uint16
	CurrentPlaymode uint8
	CurrentStatus   uint8
	StatusText      string
}

func BanchoSendOsuUpdate(packetQueue chan BanchoPacket, user database.UserStats, status OsuStatus) {
	buf := new(bytes.Buffer)

	//Write Data
	binary.Write(buf, binary.LittleEndian, int32(user.UserID))

	binary.Write(buf, binary.LittleEndian, status.CurrentStatus)
	binary.Write(buf, binary.LittleEndian, WriteBanchoString(status.StatusText))
	binary.Write(buf, binary.LittleEndian, WriteBanchoString(status.BeatmapChecksum))
	binary.Write(buf, binary.LittleEndian, status.CurrentMods)
	binary.Write(buf, binary.LittleEndian, status.CurrentPlaymode)
	binary.Write(buf, binary.LittleEndian, status.BeatmapId)

	binary.Write(buf, binary.LittleEndian, int64(user.RankedScore))
	binary.Write(buf, binary.LittleEndian, user.Accuracy)
	binary.Write(buf, binary.LittleEndian, int32(user.Playcount))
	binary.Write(buf, binary.LittleEndian, int64(user.TotalScore))
	binary.Write(buf, binary.LittleEndian, int32(1)) //TODO: rank

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
