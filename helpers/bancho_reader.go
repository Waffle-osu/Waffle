package helpers

import "io"

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
