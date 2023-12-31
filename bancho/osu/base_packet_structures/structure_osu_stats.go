package base_packet_structures

type OsuStats struct {
	UserId      int32
	Status      StatusUpdate
	RankedScore int64
	Accuracy    float32
	Playcount   int32
	TotalScore  int64
	Rank        int32
}
