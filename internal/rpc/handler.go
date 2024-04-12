package rpc

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

// handleWebSocket handles incoming WebSocket connections
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
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
