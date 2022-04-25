package packets

import (
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

	message.Sender = string(ReadBanchoString(reader))
	message.Message = string(ReadBanchoString(reader))
	message.Target = string(ReadBanchoString(reader))

	return message
}

func (message *Message) WriteMessage(writer io.Writer) {
	binary.Write(writer, binary.LittleEndian, WriteBanchoString(message.Sender))
	binary.Write(writer, binary.LittleEndian, WriteBanchoString(message.Message))
	binary.Write(writer, binary.LittleEndian, WriteBanchoString(message.Target))
}
