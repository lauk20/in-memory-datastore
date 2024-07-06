package main

import (
	"sync"
)

// map of command to function handler
var handlers = map[string]func([]Value) Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
}

// map from string key to string value
var sets = map[string]string{}

// mutex for sets
var setsMutex = sync.RWMutex{}

// assign a value to a key in the sets map
// returns Value with message
func set(args []Value) Value {
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
func get(args []Value) Value {
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

// handler for PING command
// returns PONG if no arguments given
func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{valueType: "string", str: "PONG"}
	}

	return Value{valueType: "string", str: args[0].bulk}
}
