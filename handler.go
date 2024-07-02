package main

// map of command to function handler
var handlers = map[string]func([]Value) Value{
	"PING": ping,
}

// handler for PING command
// returns PONG if no arguments given
func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{valueType: "string", str: "PONG"}
	}

	return Value{valueType: "string", str: args[0].bulk}
}
