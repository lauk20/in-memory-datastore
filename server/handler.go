package main

import (
	"sync"
)

// server type to store de/serializers
// used to maintain a connection between client and server smoothly
type Server struct {
	deserializer *Deserializer
	serializer   *Serializer
}

// map of command to function handler
var handlers = map[string]func([]Value, *Server) Value{
	"PING":      ping,
	"SET":       set,
	"GET":       get,
	"SUBSCRIBE": subscribe,
	"PUBLISH":   publish,
}

// pub/sub message broker
var broker = NewBroker()

// map from string key to string value
var sets = map[string]string{}

// mutex for sets
var setsMutex = sync.RWMutex{}

// assign a value to a key in the sets map
// returns Value with message
func set(args []Value, s *Server) Value {
	if len(args) != 2 {
		return Value{valueType: "error", str: "Invalid arg count for set"}
	}

	key := args[0].bulk
	value := args[1].bulk

	setsMutex.Lock()
	sets[key] = value
	setsMutex.Unlock()

	return Value{valueType: "string", str: "OK"}
}

// get a value using key in sets map
// returns Value with message
func get(args []Value, s *Server) Value {
	if len(args) != 1 {
		return Value{valueType: "error", str: "Invaid arg count for get"}
	}

	key := args[0].bulk

	setsMutex.RLock()
	value, found := sets[key]
	setsMutex.RUnlock()

	if !found {
		return Value{valueType: "null"}
	}

	return Value{valueType: "bulk", bulk: value}
}

func subscribe(args []Value, s *Server) Value {
	if len(args) != 1 {
		return Value{valueType: "error", str: "Invalid arg count for subscribe"}
	}

	topic := args[0].bulk
	channel := broker.Subscribe(topic)

	for {
		result := <-channel
		s.serializer.Write(result)
	}

	return Value{valueType: "string", str: "subscribe done"}
}

func publish(args []Value, s *Server) Value {
	if len(args) != 2 {
		return Value{valueType: "error", str: "Invalid arg count for publish"}
	}

	topic := args[0].bulk
	msg := args[1]

	broker.Publish(topic, msg)

	return Value{valueType: "string", str: "Published message"}
}

// handler for PING command
// returns PONG if no arguments given
func ping(args []Value, s *Server) Value {
	if len(args) == 0 {
		return Value{valueType: "string", str: "PONG"}
	}

	return Value{valueType: "string", str: args[0].bulk}
}
