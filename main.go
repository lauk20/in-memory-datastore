package main

import (
	"fmt"
	"net"
	"strings"
)

// server event loop
func serverLoop(connection net.Conn) {
	// server listen loop
	for {
		deserializer := NewDeserializer(connection)
		value, err := deserializer.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		// requests should be of type array "set key 1" is an array
		if value.valueType != "array" {
			fmt.Println("Invalid request")
			continue
		}

		// the request should not be empty
		if len(value.array) == 0 {
			fmt.Println("Invalid request")
			continue
		}

		// case insensitive commands
		command := strings.ToUpper(value.array[0].bulk)
		// arguments
		args := value.array[1:]

		// get Serializer for this connection
		serializer := NewSerializer(connection)

		// get the handler for this command
		handlerFunction, found := handlers[command]
		if !found {
			fmt.Println("Invalid command", command)
			serializer.Write(Value{valueType: "string", str: ""})
			continue
		}

		// run handler with args
		result := handlerFunction(args)
		// respond with result
		serializer.Write(result)
	}
}

func main() {
	fmt.Println("Listening on port 6379")

	// start listening on port 6379
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		fmt.Println("Listening...")
		// block for connection
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		// close connection once function exits
		defer connection.Close()

		// start goroutine
		go serverLoop(connection)
	}
}
