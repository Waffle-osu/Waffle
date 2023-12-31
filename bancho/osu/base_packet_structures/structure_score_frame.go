package base_packet_structures

type ScoreFrame struct {
	Time         int32
	Id           uint8
	Count300     uint16
	Count100     uint16
	Count50      uint16
	CountGeki    uint16
	CountKatu    uint16
	CountMiss    uint16
	TotalScore   int32
	MaxCombo     uint16
	CurrentCombo uint16
	Perfect      bool
	CurrentHp    uint8
	TagByte      uint8
}
