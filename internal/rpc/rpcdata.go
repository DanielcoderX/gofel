// Package rpc implements the WebSocket RPC server for Gofel.
package rpc

import (
	"sync"

	"github.com/DanielcoderX/gofel/internal/utils"
	"github.com/gorilla/websocket"
)

// RPCMap defines the structure for function mappings.
type RPCMap map[string]func(*websocket.Conn, interface{}) error

// RPCFunctions stores the mapping of function names to function calls.
var RPCFunctions = struct {
	sync.RWMutex
	m RPCMap
}{m: make(RPCMap)}

// On registers a callback function that will be executed when called via RPC.
func On(name string, callback func(conn *websocket.Conn, data interface{})) {
	RPCFunctions.Lock()
	defer RPCFunctions.Unlock()
	RPCFunctions.m[name] = func(conn *websocket.Conn, data interface{}) error {
		callback(conn, data)
		return nil
	}
}

// callRPCFunction tries to call a registered RPC function by name with provided data.
func callRPCFunction(funcName string, conn *websocket.Conn, data interface{}, verbose bool) {
	RPCFunctions.RLock()
	function, exists := RPCFunctions.m[funcName]
	RPCFunctions.RUnlock()
	if exists {
		go function(conn, data) // Run in a goroutine to handle concurrently
	} else {
		conn.WriteMessage(websocket.TextMessage, []byte("Function "+funcName+" not found"))
		utils.LogVerbose(verbose, "Function %s not found", funcName)
	}
}
