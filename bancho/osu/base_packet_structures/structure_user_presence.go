package base_packet_structures

type UserPresence struct {
	UserId          int32
	Username        string
	AvatarExtension uint8
	Timezone        uint8
	Country         uint8
	City            string
	Permissions     uint8
	Longitude       float32
	Latitude        float32
	Rank            int32
}
