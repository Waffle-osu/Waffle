package packets

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"
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

func isNil[T any](t T) bool {
	v := reflect.ValueOf(t)
	kind := v.Kind()
	// Must be one of these types to be nillable
	return (kind == reflect.Ptr ||
		kind == reflect.Interface ||
		kind == reflect.Slice ||
		kind == reflect.Map ||
		kind == reflect.Chan ||
		kind == reflect.Func) &&
		v.IsNil()
}

func Send[T any](packetId uint16, serializable T) []byte {
	packetBuffer := new(bytes.Buffer)

	//Packet ID and Compression
	binary.Write(packetBuffer, binary.LittleEndian, packetId)
	binary.Write(packetBuffer, binary.LittleEndian, uint8(0))

	if isNil(serializable) {
		binary.Write(packetBuffer, binary.LittleEndian, 0)
	} else {
		packetBytes := ReflectionWrite(serializable)
		packetLength := len(packetBytes)

		binary.Write(packetBuffer, binary.LittleEndian, int32(packetLength))
		binary.Write(packetBuffer, binary.LittleEndian, packetBytes)
	}

	return packetBuffer.Bytes()
}

func SendEmpty(packetId uint16) []byte {
	return Send[struct{}](packetId, struct{}{})
}

func SendBytes(packetId uint16, sendBytes []byte) []byte {
	packetLength := len(sendBytes)

	packetBuffer := new(bytes.Buffer)

	binary.Write(packetBuffer, binary.LittleEndian, packetId)
	binary.Write(packetBuffer, binary.LittleEndian, uint8(0))
	binary.Write(packetBuffer, binary.LittleEndian, int32(packetLength))

	binary.Write(packetBuffer, binary.LittleEndian, sendBytes)

	return packetBuffer.Bytes()
}
