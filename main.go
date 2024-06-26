package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	// start listening on port 6379
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	// block for connection
	connection, err := listener.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	// close connection once function exits
	defer connection.Close()

	// server listen loop
	for {
		// buffer of 1024 bytes
		buffer := make([]byte, 1024)

		// read data into buffer
		_, err = connection.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("client sent invalid data: ", err.Error())
		}

		// respond with "+OK\r\n"
		connection.Write([]byte("+OK\r\n"))
	}
}
