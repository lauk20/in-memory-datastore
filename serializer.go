package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// definitions for each symbol token
const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

// type to hold the data for each request
type Value struct {
	valueType string  // string indicating the type of request
	str       string  // string for basic strings
	num       int     // int for numbers
	bulk      string  // string for bulk strings
	array     []Value // []Value holds Values for arrays
}

// type for Deserializing the protocol's input
type Deserializer struct {
	reader *bufio.Reader
}

type Serializer struct {
	writer io.Writer
}

// create new Deserializer
func NewDeserializer(reader io.Reader) *Deserializer {
	return &Deserializer{reader: bufio.NewReader(reader)}
}

// create new Serializer
func NewSerializer(writer io.Writer) *Serializer {
	return &Serializer{writer: writer}
}

// read input line until CRLF
// returns the line, bytes read, and error if any
func (d *Deserializer) readLine() (line []byte, n int, err error) {
	for {
		readByte, err := d.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, readByte)
		// check for CRLF
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	// return the line, removing the \r
	return line[:len(line)-2], n, nil
}

// read input for integer inputs
// returns the integer read
func (d *Deserializer) readInteger() (x int, n int, err error) {
	line, n, err := d.readLine()
	if err != nil {
		return 0, 0, err
	}
	num, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(num), n, nil
}

// read array input
// returns the Value
func (d *Deserializer) readArray() (Value, error) {
	// create Value data
	v := Value{}
	v.valueType = "array"

	// get array length
	length, _, err := d.readInteger()
	if err != nil {
		return v, err
	}

	// init array as Value field
	v.array = make([]Value, 0)
	for i := 0; i < length; i++ {
		// read the Value in the array input
		val, err := d.Read()
		if err != nil {
			return v, err
		}

		// append Value to array
		v.array = append(v.array, val)
	}

	// return value
	return v, nil
}

// read bulk strings
// returns the Value structure
func (d *Deserializer) readBulk() (Value, error) {
	// create Value Data
	v := Value{}
	v.valueType = "bulk"

	// get bulk length
	length, _, err := d.readInteger()
	if err != nil {
		return v, err
	}

	// create []byte for bulk
	bulk := make([]byte, length)
	// read the bulk string into bulk
	d.reader.Read(bulk)
	// consume the CRLF
	d.readLine()

	// assign bulk Value
	v.bulk = string(bulk)

	return v, nil
}

// main Read function to read input
// returns the Read value
func (d *Deserializer) Read() (Value, error) {
	inputType, err := d.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch inputType {
	case ARRAY:
		return d.readArray()
	case BULK:
		return d.readBulk()
	default:
		fmt.Printf("Invalid type: %v", string(inputType))
		return Value{}, nil
	}
}

// serialize an Array value
// return []byte after serialized
func (v Value) serializeArray() []byte {
	var result = []byte{}
	result = append(result, ARRAY)
	result = append(result, strconv.Itoa(len(v.array))...)
	result = append(result, '\r', '\n')
	for i := 0; i < len(v.array); i++ {
		result = append(result, v.array[i].Serialize()...)
	}

	return result
}

// serialize a Bulk value
// return []byte after serialized
func (v Value) serializeBulk() []byte {
	var result []byte
	result = append(result, BULK)
	result = append(result, strconv.Itoa(len(v.bulk))...)
	result = append(result, '\r', '\n')
	result = append(result, v.bulk...)
	result = append(result, '\r', '\n')

	return result
}

// serialize a String value
// return []byte after serialized
func (v Value) serializeString() []byte {
	var result []byte
	result = append(result, STRING)
	result = append(result, v.str...)
	result = append(result, '\r', '\n')

	return result
}

// serialize a Null value
// return []byte after serialized
func (v Value) serializeNull() []byte {
	return []byte("$-1\r\n")
}

// serialize an Error value
// return []byte after serialized
func (v Value) serializeError() []byte {
	var result []byte
	result = append(result, ERROR)
	result = append(result, v.str...)
	result = append(result, '\r', '\n')

	return result
}

// serializer an Integer value
// return []byte after serialized
func (v Value) serializeInteger() []byte {
	var result []byte
	result = append(result, INTEGER)
	result = append(result, strconv.Itoa(v.num)...)
	result = append(result, '\r', '\n')

	return result
}

// function to convert Value to serialized bytes array
// returns serialized []byte
func (v Value) Serialize() []byte {
	switch v.valueType {
	case "array":
		return v.serializeArray()
	case "bulk":
		return v.serializeBulk()
	case "string":
		return v.serializeString()
	case "null":
		return v.serializeNull()
	case "error":
		return v.serializeError()
	case "integer":
		return v.serializeInteger()
	default:
		return []byte{}
	}
}

// Serialize and write a value using Serializer's writer
// return error if not success
func (s *Serializer) Write(v Value) error {
	var result = v.Serialize()
	_, err := s.writer.Write(result)
	if err != nil {
		return err
	}

	return nil
}
