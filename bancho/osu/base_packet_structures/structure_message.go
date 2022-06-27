package base_packet_structures

import (
	"Waffle/helpers"
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

	message.Sender = string(helpers.ReadBanchoString(reader))
	message.Message = string(helpers.ReadBanchoString(reader))
	message.Target = string(helpers.ReadBanchoString(reader))

	return message
}

func (message *Message) WriteMessage(writer io.Writer) {
	binary.Write(writer, binary.LittleEndian, helpers.WriteBanchoString(message.Sender))
	binary.Write(writer, binary.LittleEndian, helpers.WriteBanchoString(message.Message))
	binary.Write(writer, binary.LittleEndian, helpers.WriteBanchoString(message.Target))
}
