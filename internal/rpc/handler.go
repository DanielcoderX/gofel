// Package rpc implements the WebSocket RPC server for gofel.
package rpc

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// upgrader is the WebSocket upgrader with some custom settings.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

// HandleWebSocket handles incoming WebSocket connections.
//
// It upgrades the HTTP connection to the WebSocket protocol,
// reads JSON-RPC requests from the client, unmarshals the requests,
// calls the corresponding function, and sends the response back to the client.
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade websocket: %s", err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Println("WebSocket closed normally")
			} else {
				log.Printf("Read error: %v", err)
			}
			break
		}

		var rpcRequest map[string]interface{}
		err = json.Unmarshal(msg, &rpcRequest)
		if err != nil {
			log.Println("JSON unmarshal error:", err)
			continue
		}

		if funcName, ok := rpcRequest["function"].(string); ok {
			callRPCFunction(funcName, conn, rpcRequest["data"])
		}
	}
}

