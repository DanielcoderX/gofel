package api

import (
	"context"
	"net/http"

	"github.com/DanielcoderX/gofel/internal/rpc"
	"github.com/DanielcoderX/gofel/internal/utils"
	"github.com/DanielcoderX/gofel/pkg/config"

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

func (s *Server) On(name string, callback func(*websocket.Conn, interface{})) {
	rpc.On(name, callback)
}

// Start starts the server and listens for incoming connections.
func (s *Server) Start(ctx context.Context) error {
	http.HandleFunc(s.config.Path, func(w http.ResponseWriter, r *http.Request) {
		rpc.HandleWebSocket(ctx, w, r, s.config.Verbose)
	})

	utils.LogVerbose(s.config.Verbose, "RPC starting on %s...", s.config.Port)

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
