package base_packet_structures

import (
	"Waffle/helpers/serialization"
	"encoding/binary"
	"io"
)

type Message struct {
	Sender  string
	Message string
	Target  string
}

func ReadMessage(reader io.Reader) Message {
	message := Message{}

	message.Sender = string(serialization.ReadBanchoString(reader))
	message.Message = string(serialization.ReadBanchoString(reader))
	message.Target = string(serialization.ReadBanchoString(reader))

	return message
}

func (message Message) Write(writer io.Writer) {
	binary.Write(writer, binary.LittleEndian, serialization.WriteBanchoString(message.Sender))
	binary.Write(writer, binary.LittleEndian, serialization.WriteBanchoString(message.Message))
	binary.Write(writer, binary.LittleEndian, serialization.WriteBanchoString(message.Target))
}
