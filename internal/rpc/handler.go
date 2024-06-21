package rpc

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/DanielcoderX/gofel/internal/utils"
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
func HandleWebSocket(ctx context.Context, w http.ResponseWriter, r *http.Request, verbose bool, format string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.LogVerbose(verbose, "Failed to upgrade websocket: %s", err)
		return
	}
	defer conn.Close()

	wrappedConn := wsconn.NewConnectionWrapper(conn)

	utils.LogVerbose(verbose, "WebSocket connection established")

	for {
		select {
		case <-ctx.Done():
			utils.LogVerbose(verbose, "Context canceled, closing WebSocket connection")
			return
		default:
			_, msg, err := wrappedConn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					utils.LogVerbose(verbose, "WebSocket closed normally")
				} else {
					if strings.Contains(err.Error(), "close") {
						utils.LogVerbose(verbose, "WebSocket connection closed")
					} else {
						utils.LogVerbose(verbose, "Read error: %v", err)
					}
				}
				return // Exit the loop on any read error
			}

			var rpcRequest map[string]interface{}
			if format == "msgpack" {
				err = msgpack.Unmarshal(msg, &rpcRequest)
			} else {
				err = json.Unmarshal(msg, &rpcRequest)
			}
			if err != nil {
				utils.LogVerbose(verbose, "%s unmarshal error: %v", format, err)
				continue
			}

			if funcName, ok := rpcRequest["Function"].(string); ok {
				go callRPCFunction(funcName, wrappedConn, rpcRequest["Data"], verbose)
			}
		}
	}
}
