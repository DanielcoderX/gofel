package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"sync"

	"github.com/DanielcoderX/gofel/internal/utils"
	"github.com/DanielcoderX/gofel/pkg/config"
	"github.com/DanielcoderX/gofel/pkg/wsconn"
	"github.com/gorilla/websocket"
	"github.com/vmihailenco/msgpack/v5"
)

// RPCMap defines the structure for function mappings.
type RPCMap map[string]func(*wsconn.ConnectionWrapper, interface{}) error

// RPCFunctions stores the mapping of function names to function calls.
var RPCFunctions = struct {
	sync.RWMutex
	m RPCMap
}{m: make(RPCMap)}

// On registers a callback function that will be executed when called via RPC.
func On(name string, callback func(conn *wsconn.ConnectionWrapper, data interface{})) {
	RPCFunctions.Lock()
	defer RPCFunctions.Unlock()
	RPCFunctions.m[name] = func(conn *wsconn.ConnectionWrapper, data interface{}) error {
		callback(conn, data)
		return nil
	}
}

// MarshalData encodes data into the specified format (either "json" or "msgpack").
func MarshalData(data interface{}) ([]byte, error) {
	switch config.GlobalConfig.Format {
	case "json":
		var buff bytes.Buffer
		enc := json.NewEncoder(&buff)
		enc.SetEscapeHTML(false)
		err := enc.Encode(data)
		if err != nil {
			return nil, err
		}
		return buff.Bytes(),err
	case "msgpack":
		return msgpack.Marshal(data)
	default:
		return nil, errors.New("unsupported format: " + config.GlobalConfig.Format)
	}
}

// callRPCFunction tries to call a registered RPC function by name with provided data.
func callRPCFunction(funcName string, conn *wsconn.ConnectionWrapper, data interface{}, verbose bool) {
	RPCFunctions.RLock()
	function, exists := RPCFunctions.m[funcName]
	RPCFunctions.RUnlock()
	if exists {
		go func() {
			err := function(conn, data) // Run in a goroutine to handle concurrently
			if err != nil {
				conn.SendMessage(websocket.TextMessage, []byte("Error: "+err.Error()))
				utils.LogVerbose(verbose, "Error calling function %s: %v", funcName, err)
			}
		}()
	} else {
		conn.SendMessage(websocket.TextMessage, []byte("Function "+funcName+" not found"))
		utils.LogVerbose(verbose, "Function %s not found", funcName)
	}
}
