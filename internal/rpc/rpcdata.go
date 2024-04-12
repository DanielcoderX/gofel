// Package rpc implements the WebSocket RPC server for gofel.
package rpc

import (
	"log"

	"github.com/gorilla/websocket"
)

// RPCMap type defines the structure for function mappings
//
// RPCMap is a map of function names to functions that can be called via RPC
type RPCMap map[string]func(*websocket.Conn, interface{}) error

// RPCFunctions stores the mapping of function names to function calls
//
// RPCFunctions is a variable that holds the mapping of function names to
// functions that can be called via RPC. It is a RPCMap type.
var RPCFunctions RPCMap

// RegisterFunction allows users to register functions that can be called via RPC
//
// RegisterFunction allows users to register functions that can be called via RPC.
// It takes two parameters, a string name, and a function that can be called
// via RPC. The function takes two parameters, a *websocket.Conn and an interface{}.
//
// If the RPCFunctions variable is nil, it is initialized as a new RPCMap. The
// function is then added to the RPCMap under the key specified by the name
// parameter.
func RegisterFunction(name string, function func(*websocket.Conn, interface{}) error) {
	if RPCFunctions == nil {
		RPCFunctions = make(RPCMap)
	}
	RPCFunctions[name] = function
}

// CallRPCFunction tries to call a registered RPC function by name with provided data
//
// CallRPCFunction tries to call a registered RPC function by name with provided
// data. If the function is found in the RPCFunctions map, it is called with the
// provided data. If the function is not found, an error is sent to the client
// and logged.
func callRPCFunction(funcName string, conn *websocket.Conn, data interface{}) {
	if function, exists := RPCFunctions[funcName]; exists {
		go function(conn, data) // Run in a goroutine to handle concurrently
	} else {
		conn.WriteMessage(websocket.TextMessage, []byte("Function "+funcName+" not found"))
		log.Printf("Function %s not found", funcName)
	}
}

