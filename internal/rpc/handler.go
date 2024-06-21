package rpc

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/DanielcoderX/gofel/internal/utils"
	"github.com/DanielcoderX/gofel/pkg/config"
	"github.com/DanielcoderX/gofel/pkg/wsconn"
	"github.com/gorilla/websocket"
	"github.com/vmihailenco/msgpack/v5"
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
func HandleWebSocket(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.LogVerbose(config.GlobalConfig.Verbose, "Failed to upgrade websocket: %s", err)
		return
	}
	defer conn.Close()

	// Get connection from the pool
	wrappedConn, err := wsconn.GlobalPool.GetConnection(conn)
	if err != nil {
		utils.LogVerbose(config.GlobalConfig.Verbose, "Failed to get connection from pool: %s", err)
		return
	}
	defer wsconn.GlobalPool.ReleaseConnection(wrappedConn)

	utils.LogVerbose(config.GlobalConfig.Verbose, "WebSocket connection established")

	for {
		select {
		case <-ctx.Done():
			utils.LogVerbose(config.GlobalConfig.Verbose, "Context canceled, closing WebSocket connection")
			return
		default:
			_, msg, err := wrappedConn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					utils.LogVerbose(config.GlobalConfig.Verbose, "WebSocket closed normally")
				} else {
					if strings.Contains(err.Error(), "close") {
						utils.LogVerbose(config.GlobalConfig.Verbose, "WebSocket connection closed")
					} else {
						utils.LogVerbose(config.GlobalConfig.Verbose, "Read error: %v", err)
					}
				}
				return // Exit the loop on any read error
			}
			var rpcRequest map[string]interface{}
			if config.GlobalConfig.Format == "msgpack" {
				err = msgpack.Unmarshal(msg, &rpcRequest)
			} else {
				err = json.Unmarshal(msg, &rpcRequest)
			}
			if err != nil {
				utils.LogVerbose(config.GlobalConfig.Verbose, "%s unmarshal error: %v", config.GlobalConfig.Format, err)
				continue
			}
			if funcName, ok := rpcRequest["function"].(string); ok {
				go callRPCFunction(funcName, wrappedConn, rpcRequest["data"], config.GlobalConfig.Verbose)
			}
		}
	}
}
