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

// create new Serializer
func NewDeserializer(reader io.Reader) *Deserializer {
	return &Deserializer{reader: bufio.NewReader(reader)}
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
