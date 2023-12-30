package serialization

import (
	"bytes"
	"encoding/binary"
	"io"
)

type BanchoSerializable interface {
	Write(writer io.Writer)
}

func WriteBanchoString(value string) []byte {
	if value == "" {
		return []byte{0}
	}

	var length int
	var i int = len(value)
	var ulebBytes []byte

	if i == 0 {
		ulebBytes = []byte{0}
	}

	for i > 0 {
		ulebBytes = append(ulebBytes, 0)
		ulebBytes[length] = byte(i & 0x7F)
		i >>= 7
		if i != 0 {
			ulebBytes[length] |= 0x80
		}
		length++
	}

	returnBytes := []byte{11}
	returnBytes = append(returnBytes, ulebBytes...)
	returnBytes = append(returnBytes, []byte(value)...)

	return returnBytes
}

func ReadBanchoString(reader io.Reader) []byte {
	bytes := make([]byte, 1)

	reader.Read(bytes)

	if bytes[0] != 11 {
		return []byte{}
	}

	var shift uint
	var lastByte byte
	var total int

	for {
		b := make([]byte, 1)
		reader.Read(b)
		lastByte = b[0]
		total |= (int(lastByte&0x7F) << shift)
		if lastByte&0x80 == 0 {
			break
		}
		shift += 7
	}

	bytes = make([]byte, total)

	reader.Read(bytes)

	return bytes
}

func SendSerializable(packetId uint16, serializable any) []byte {
	buf := new(bytes.Buffer)

	ReflectionWrite(serializable)

	packetBytes := buf.Bytes()
	packetLength := len(packetBytes)

	packetBuffer := new(bytes.Buffer)

	binary.Write(packetBuffer, binary.LittleEndian, packetId)
	binary.Write(packetBuffer, binary.LittleEndian, uint8(0))
	binary.Write(packetBuffer, binary.LittleEndian, int32(packetLength))

	binary.Write(packetBuffer, binary.LittleEndian, packetBytes)

	return packetBuffer.Bytes()
}

func SendSerializableString(packetId uint16, str string) []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, WriteBanchoString(str))

	packetBytes := buf.Bytes()
	packetLength := len(packetBytes)

	packetBuffer := new(bytes.Buffer)

	binary.Write(packetBuffer, binary.LittleEndian, packetId)
	binary.Write(packetBuffer, binary.LittleEndian, uint8(0))
	binary.Write(packetBuffer, binary.LittleEndian, int32(packetLength))

	binary.Write(packetBuffer, binary.LittleEndian, packetBytes)

	return packetBuffer.Bytes()
}

func SendSerializableBytes(packetId uint16, sendBytes []byte) []byte {
	packetLength := len(sendBytes)

	packetBuffer := new(bytes.Buffer)

	binary.Write(packetBuffer, binary.LittleEndian, packetId)
	binary.Write(packetBuffer, binary.LittleEndian, uint8(0))
	binary.Write(packetBuffer, binary.LittleEndian, int32(packetLength))

	binary.Write(packetBuffer, binary.LittleEndian, sendBytes)

	return packetBuffer.Bytes()
}

func SendSerializableInt(packetId uint16, integer int32) []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, packetId)
	binary.Write(buf, binary.LittleEndian, uint8(0))
	binary.Write(buf, binary.LittleEndian, int32(4))
	binary.Write(buf, binary.LittleEndian, integer)

	return buf.Bytes()
}

func SendEmptySerializable(packetId uint16) []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, packetId)
	binary.Write(buf, binary.LittleEndian, uint8(0))
	binary.Write(buf, binary.LittleEndian, int32(0))

	return buf.Bytes()
}
