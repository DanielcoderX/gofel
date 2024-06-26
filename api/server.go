package api

import (
	"context"
	"net/http"

	"github.com/DanielcoderX/gofel/internal/rpc"
	"github.com/DanielcoderX/gofel/internal/utils"
	"github.com/DanielcoderX/gofel/pkg/config"
	"github.com/DanielcoderX/gofel/pkg/wsconn"
	"github.com/gorilla/websocket"
)

type Server struct {
	config     *config.Config
	httpServer *http.Server // Store a reference to the HTTP server
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

// On registers a callback function for an RPC call.
func (s *Server) On(name string, callback func(*wsconn.ConnectionWrapper, interface{})) {
	rpc.On(name, callback)
}

// SendResponse marshals the data and sends it over the WebSocket connection
// using the appropriate message type based on the configured format.
func SendResponse(conn *wsconn.ConnectionWrapper, data interface{}) error {
	// Marshal the data based on the format specified in the config.
	response, err := rpc.MarshalData(data)
	if err != nil {
		return err
	}

	// Determine the message type based on the config format.
	messageType := websocket.TextMessage
	if config.GlobalConfig.Format == "msgpack" {
		messageType = websocket.BinaryMessage
	}

	// Send the marshaled data as a WebSocket message.
	return conn.SendMessage(messageType, response)
}

// Start starts the server and listens for incoming connections.
func (s *Server) Start(ctx context.Context) error {
	// Globalize the config
	config.GlobalConfig = s.config
	// Initialize the global connection pool with a capacity from config
	wsconn.InitConnectionPool(s.config.ConnectionPoolCapacity)
	http.HandleFunc(s.config.Path, func(w http.ResponseWriter, r *http.Request) {
		rpc.HandleWebSocket(ctx, w, r)
	})

	utils.LogVerbose(s.config.Verbose, "RPC starting on %s at %s", s.config.Port, s.config.Path)

	s.httpServer = &http.Server{Addr: ":" + s.config.Port}

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop gracefully shuts down the server.
func (s *Server) Stop() error {
	utils.LogVerbose(s.config.Verbose, "Shutting down server gracefully...")

	// Shut down the HTTP server gracefully
	if err := s.httpServer.Shutdown(context.TODO()); err != nil {
		return err
	}

	utils.LogVerbose(s.config.Verbose, "Server shut down gracefully.")
	return nil
}
