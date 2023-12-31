package base_packet_structures

type SpectatorFrame struct {
	ButtonState           uint8
	ButtonStateCompatByte uint8
	MouseX                float32
	MouseY                float32
	Time                  int32
}
