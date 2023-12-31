package serialization

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"
	"strconv"
)

// Reads a single value
func writeSingle(writer io.Writer, indirect reflect.Value, elements reflect.Type, field reflect.StructField, fieldValue reflect.Value) {
	//Dereferenced value of the field
	v := fieldValue.Addr().Interface()

	//Each type has to be converted differently
	switch fieldValue.Kind() {
	case reflect.Bool:
		conv, _ := (v.(*bool))

		binary.Write(writer, binary.LittleEndian, *conv)
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

		asBanchoString := WriteBanchoString(*conv)

		binary.Write(writer, binary.LittleEndian, asBanchoString)
	case reflect.Slice:
		//Slices rely on a previously defined field specified in the tag
		//this gets the tag value
		lengthField := field.Tag.Get("length")

		if lengthField == "" {
			panic("No length field specified!")
		}

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
	case reflect.Array:
		//Special handling for Multi Slot statuses
		multiSpecialField := field.Tag.Get("multi")
		multiAndCondition := field.Tag.Get("multiAnd")
		var intMultiAndCond int64

		if multiAndCondition != "" {
			parsedIntCond, err := strconv.ParseInt(multiAndCondition, 10, 64)

			if err != nil {
				panic("Multi and condition not int.")
			}

			intMultiAndCond = parsedIntCond
		}

		arrayLength := fieldValue.Len()

		specialField := indirect.FieldByName(multiSpecialField)

		var slotStatus []uint8

		//Populate slotStatus if available
		if multiSpecialField != "" {
			switch specialField.Type().String() {
			case "[8]uint8":
				//SlotStatus arrays
				slotStatusParsed, ok := specialField.Addr().Elem().Interface().([8]uint8)

				if !ok {
					panic("Something's wrong with SlotStatuses")
				}

				slotStatus = slotStatusParsed[:]
			case "[16]uint8":
				//SlotStatus arrays
				slotStatusParsed, ok := specialField.Addr().Elem().Interface().([16]uint8)

				if !ok {
					panic("Something's wrong with SlotStatuses")
				}

				slotStatus = slotStatusParsed[:]
			}
		}

		for j := 0; j != arrayLength; j++ {
			//No Special multi handling
			if multiSpecialField == "" {
				indexValue := fieldValue.Index(j)

				//Reusing this function to write it
				writeSingle(writer, indexValue, indexValue.Type(), field, indexValue)
			} else {
				if (slotStatus[j] & uint8(intMultiAndCond)) > 0 {
					indexValue := fieldValue.Index(j)

					//Reusing this function to write it
					writeSingle(writer, indexValue, indexValue.Type(), field, indexValue)
				}
			}
		}
	}
}

func writeInternal(writer io.Writer, indirect reflect.Value, elements reflect.Type) {
	for i := 0; i != elements.NumField(); i++ {
		//General field type information, used for getting the type and the Tag
		field := elements.Field(i)
		//Actual value of the field
		fieldValue := indirect.Field(i)

		switch field.Type.Kind() {
		case reflect.Struct:
			//Recursive call so we can handle structs in structs
			writeInternal(writer, indirect.Field(i), field.Type)
		default:
			writeSingle(writer, indirect, elements, field, fieldValue)
		}
	}
}

func ReflectionWrite[T any](packetStruct T) []byte {
	buffer := new(bytes.Buffer)
	//Holds all the values
	indirect := reflect.Indirect(reflect.ValueOf(&packetStruct))
	//Holds the type information, needed for tags
	elements := indirect.Type()

	writeInternal(buffer, indirect, elements)

	return buffer.Bytes()
}

func readSingle(reader io.Reader, val reflect.Value, structField reflect.StructField, field reflect.Value) {
	kind := field.Kind()

	switch kind {
	case reflect.Pointer:
		elemed := field.Elem()

		readSingle(reader, elemed, structField, elemed)
	case reflect.Bool:
		var value bool

		binary.Read(reader, binary.LittleEndian, &value)

		field.Set(reflect.ValueOf(value))
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
		string := string(ReadBanchoString(reader))

		field.Set(reflect.ValueOf(string))
	case reflect.Struct:
		readInternal(reader, field)
	case reflect.Slice:
		//Slices rely on a previously defined field specified in the tag
		//this gets the tag value
		lengthField := structField.Tag.Get("length")

		if lengthField == "" {
			panic("No length field specified!")
		}

		//Field information
		lengthFieldByName := val.FieldByName(lengthField)
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

		elemSlice := reflect.MakeSlice(reflect.SliceOf(field.Type().Elem()), 0, int(iterations))

		for j := uint64(0); j != iterations; j++ {
			value := reflect.New(field.Type().Elem())

			//String hotfix
			if value.Type().String() == "*string" {
				value = value.Elem()

				readSingle(reader, value, structField, value)

				elemSlice = reflect.Append(elemSlice, value)
			} else {
				readSingle(reader, value, structField, value)

				elemSlice = reflect.Append(elemSlice, value.Elem())
			}
		}

		field.Set(elemSlice)
	case reflect.Array:
		//Special handling for Multi Slot statuses
		multiSpecialField := structField.Tag.Get("multi")
		multiAndCondition := structField.Tag.Get("multiAnd")
		var intMultiAndCond int64

		if multiAndCondition != "" {
			parsedIntCond, err := strconv.ParseInt(multiAndCondition, 10, 64)

			if err != nil {
				panic("Multi and condition not int.")
			}

			intMultiAndCond = parsedIntCond
		}

		arrayLength := field.Len()

		specialField := val.FieldByName(multiSpecialField)

		var slotStatus []uint8

		//Populate slotStatus if available
		if multiSpecialField != "" {
			switch specialField.Type().String() {
			case "[8]uint8":
				//SlotStatus arrays
				slotStatusParsed, ok := specialField.Addr().Elem().Interface().([8]uint8)

				if !ok {
					panic("Something's wrong with SlotStatuses")
				}

				slotStatus = slotStatusParsed[:]
			case "[16]uint8":
				//SlotStatus arrays
				slotStatusParsed, ok := specialField.Addr().Elem().Interface().([16]uint8)

				if !ok {
					panic("Something's wrong with SlotStatuses")
				}

				slotStatus = slotStatusParsed[:]
			}
		}

		for j := 0; j != arrayLength; j++ {
			//No Special multi handling
			if multiSpecialField == "" {
				indexValue := field.Index(j)

				//Reusing this function to write it
				readSingle(reader, indexValue, structField, indexValue)
			} else {
				if (slotStatus[j] & uint8(intMultiAndCond)) > 0 {
					indexValue := field.Index(j)

					//Reusing this function to write it
					readSingle(reader, indexValue, structField, indexValue)
				}
			}
		}
	}
}

func readInternal(reader io.Reader, val reflect.Value) {
	if val.Kind() == reflect.Interface && !val.IsNil() {
		e := val.Elem()

		if e.Kind() == reflect.Pointer && !e.IsNil() {
			val = e
		}
	}

	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}

		val = val.Elem()
	}

	elements := val.Type()

	for i := 0; i != val.NumField(); i++ {
		field := val.Field(i)
		structField := elements.Field(i)

		switch field.Kind() {
		case reflect.Struct:
			readInternal(reader, field)
		default:
			readSingle(reader, val, structField, field)
		}
	}
}

func reflectionRead(reader io.Reader, v any) {
	readInternal(reader, reflect.ValueOf(v))
}

func Read[T any](reader io.Reader) T {
	var out T

	reflectionRead(reader, &out)

	return out
}
