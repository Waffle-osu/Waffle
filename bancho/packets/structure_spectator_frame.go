package packets

import (
	"encoding/binary"
	"io"
)

type SpectatorFrame struct {
	ButtonState           uint8
	ButtonStateCompatByte uint8
	MouseX                float32
	MouseY                float32
	Time                  int32
}

func ReadSpectatorFrame(reader io.Reader) SpectatorFrame {
	frame := SpectatorFrame{}

	binary.Read(reader, binary.LittleEndian, &frame.ButtonState)
	binary.Read(reader, binary.LittleEndian, &frame.ButtonStateCompatByte)
	binary.Read(reader, binary.LittleEndian, &frame.MouseX)
	binary.Read(reader, binary.LittleEndian, &frame.MouseY)
	binary.Read(reader, binary.LittleEndian, &frame.Time)

	return frame
}

func (frame *SpectatorFrame) WriteSpectatorFrame(writer io.Writer) {
	binary.Write(writer, binary.LittleEndian, frame.ButtonState)
	binary.Write(writer, binary.LittleEndian, frame.ButtonStateCompatByte)
	binary.Write(writer, binary.LittleEndian, frame.MouseX)
	binary.Write(writer, binary.LittleEndian, frame.MouseY)
	binary.Write(writer, binary.LittleEndian, frame.Time)
}
