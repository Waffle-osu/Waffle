package packets

import (
	"Waffle/bancho/osu/base_packet_structures"
	"Waffle/helpers/serialization"
	"bytes"
	"encoding/binary"
	"io"
	"reflect"
)

func PacketTest() {
	testFrame := base_packet_structures.SpectatorFrame{
		Test: base_packet_structures.Test{
			ButtonState: 24,
		},
		ButtonStateCompatByte: 48,
		MouseX:                12,
		MouseY:                121,
		Time:                  1111,
	}

	test := Write(testFrame)

	out := base_packet_structures.SpectatorFrame{}

	Read(bytes.NewBuffer(test), &out)
}

func writeInternal(writer io.Writer, indirect reflect.Value, elements reflect.Type) {
	for i := 0; i != elements.NumField(); i++ {
		//General field type information, used for getting the type and the Tag
		field := elements.Field(i)
		//Actual value of the field
		fieldValue := indirect.Field(i)
		//Dereferenced value of the field
		v := fieldValue.Addr().Interface()

		//Each type has to be converted differently
		switch field.Type.Kind() {
		case reflect.Uint8:
			conv, _ := (v.(*uint8))

			binary.Write(writer, binary.LittleEndian, *conv)
		case reflect.Int8:
			conv, _ := (v.(*int8))

			binary.Write(writer, binary.LittleEndian, *conv)
		case reflect.Uint16:
			conv, _ := (v.(*uint16))

			binary.Write(writer, binary.LittleEndian, *conv)
		case reflect.Int16:
			conv, _ := (v.(*int16))

			binary.Write(writer, binary.LittleEndian, *conv)
		case reflect.Uint32:
			conv, _ := (v.(*uint32))

			binary.Write(writer, binary.LittleEndian, *conv)
		case reflect.Int32:
			conv, _ := (v.(*int32))

			binary.Write(writer, binary.LittleEndian, *conv)
		case reflect.Uint64:
			conv, _ := (v.(*uint64))

			binary.Write(writer, binary.LittleEndian, *conv)
		case reflect.Int64:
			conv, _ := (v.(*int64))

			binary.Write(writer, binary.LittleEndian, *conv)
		case reflect.Float32:
			conv, _ := (v.(*float32))

			binary.Write(writer, binary.LittleEndian, *conv)
		case reflect.Float64:
			conv, _ := (v.(*float64))

			binary.Write(writer, binary.LittleEndian, *conv)
		case reflect.String:
			conv, _ := (v.(*string))

			asBanchoString := serialization.WriteBanchoString(*conv)

			binary.Write(writer, binary.LittleEndian, asBanchoString)
		case reflect.Struct:
			//Recursive call so we can handle structs in structs
			writeInternal(writer, indirect.Field(i), field.Type)
		case reflect.Slice:
			//Slices rely on a previously defined field specified in the tag
			//this gets the tag value
			lengthField := field.Tag.Get("length")
			//Field information
			lengthFieldByName := indirect.FieldByName(lengthField)
			//Dereferenced value of field
			lengthValue := lengthFieldByName.Addr().Interface()

			iterations := uint64(0)

			//We need to get the proper value of the length parameter
			//So we can iterate `iterations` amount of times
			switch lengthFieldByName.Kind() {
			case reflect.Uint8:
				conv, _ := (lengthValue.(*uint8))

				iterations = uint64(*conv)
			case reflect.Int8:
				conv, _ := (lengthValue.(*int8))

				iterations = uint64(*conv)
			case reflect.Uint16:
				conv, _ := (lengthValue.(*uint16))

				iterations = uint64(*conv)
			case reflect.Int16:
				conv, _ := (lengthValue.(*int16))

				iterations = uint64(*conv)
			case reflect.Uint32:
				conv, _ := (lengthValue.(*uint32))

				iterations = uint64(*conv)
			case reflect.Int32:
				conv, _ := (lengthValue.(*int32))

				iterations = uint64(*conv)
			case reflect.Uint64:
				conv, _ := (lengthValue.(*uint64))

				iterations = uint64(*conv)
			case reflect.Int64:
				conv, _ := (lengthValue.(*int64))

				iterations = uint64(*conv)
			}

			for j := uint64(0); j != iterations; j++ {
				indexValue := fieldValue.Index(int(j))

				//Reusing this function to write it
				writeInternal(writer, indexValue, indexValue.Type())
			}
		}
	}
}

func Write[T any](packetStruct T) []byte {
	buffer := new(bytes.Buffer)
	//Holds all the values
	indirect := reflect.Indirect(reflect.ValueOf(&packetStruct))
	//Holds the type information, needed for tags
	elements := indirect.Type()

	writeInternal(buffer, indirect, elements)

	return buffer.Bytes()
}

func readInternal(reader io.Reader, indirect reflect.Value, elements reflect.Type) {
	for i := 0; i != elements.NumField(); i++ {
		field := indirect.Field(i)
		kind := field.Kind()

		switch kind {
		case reflect.Uint8:
			var value uint8

			binary.Read(reader, binary.LittleEndian, &value)

			field.Set(reflect.ValueOf(value))
		case reflect.Int8:
			var value int8

			binary.Read(reader, binary.LittleEndian, &value)

			field.Set(reflect.ValueOf(value))
		case reflect.Uint16:
			var value uint16

			binary.Read(reader, binary.LittleEndian, &value)

			field.Set(reflect.ValueOf(value))
		case reflect.Int16:
			var value int16

			binary.Read(reader, binary.LittleEndian, &value)

			field.Set(reflect.ValueOf(value))
		case reflect.Uint32:
			var value uint32

			binary.Read(reader, binary.LittleEndian, &value)

			field.Set(reflect.ValueOf(value))
		case reflect.Int32:
			var value int32

			binary.Read(reader, binary.LittleEndian, &value)

			field.Set(reflect.ValueOf(value))
		case reflect.Uint64:
			var value uint64

			binary.Read(reader, binary.LittleEndian, &value)

			field.Set(reflect.ValueOf(value))
		case reflect.Int64:
			var value int64

			binary.Read(reader, binary.LittleEndian, &value)

			field.Set(reflect.ValueOf(value))
		case reflect.Float32:
			var value float32

			binary.Read(reader, binary.LittleEndian, &value)

			field.Set(reflect.ValueOf(value))
		case reflect.Float64:
			var value float64

			binary.Read(reader, binary.LittleEndian, &value)

			field.Set(reflect.ValueOf(value))
		case reflect.String:
			string := string(serialization.ReadBanchoString(reader))

			field.Set(reflect.ValueOf(string))
		case reflect.Struct:
			readInternal(reader, indirect.Field(i), field.Type())
		case reflect.Slice:
		}
	}
}

func Read(reader io.Reader, v any) {
	//Creates a new T
	typeOf := reflect.TypeOf(reflect.ValueOf(v))

	readInternal(reader, reflect.ValueOf(v), typeOf)
}
