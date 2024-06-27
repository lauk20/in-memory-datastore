package main

import (
	"bufio"
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
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
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
