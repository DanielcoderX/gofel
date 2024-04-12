package rpc

import (
	"log"

	"github.com/gorilla/websocket"
)

// RPCMap type defines the structure for function mappings
type RPCMap map[string]func(*websocket.Conn, interface{}) error

// RPCFunctions stores the mapping of function names to function calls
var RPCFunctions RPCMap

// RegisterFunction allows users to register functions that can be called via RPC
func RegisterFunction(name string, function func(*websocket.Conn, interface{}) error) {
	if RPCFunctions == nil {
		RPCFunctions = make(RPCMap)
	}
	RPCFunctions[name] = function
}

// CallRPCFunction tries to call a registered RPC function by name with provided data
func callRPCFunction(funcName string, conn *websocket.Conn, data interface{}) {
	if function, exists := RPCFunctions[funcName]; exists {
		go function(conn, data) // Run in a goroutine to handle concurrently
	} else {
		conn.WriteMessage(websocket.TextMessage, []byte("Function "+funcName+" not found"))
		log.Printf("Function %s not found", funcName)
	}
}
