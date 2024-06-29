package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Listening on port 6379")

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
		deserializer := NewDeserializer(connection)
		_, err := deserializer.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		serializer := NewSerializer(connection)
		serializer.Write(Value{valueType: "string", str: "OK"})
	}
}
