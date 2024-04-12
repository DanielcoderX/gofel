package api

import (
	"github.com/DanielcoderX/gofel/internal/rpc"
	"github.com/DanielcoderX/gofel/pkg/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
)

type Server struct {
	config *config.Config
	httpServer *http.Server // Store a reference to the HTTP server
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

func (s *Server) RegisterFunction(name string, function func(*websocket.Conn, interface{}) error) {
	rpc.RegisterFunction(name, function)
}

func (s *Server) Start() error {
	http.HandleFunc(s.config.Path, rpc.HandleWebSocket)
	log.Printf("RPC starting on %s...", s.config.Port)

	s.httpServer = &http.Server{Addr: ":" + s.config.Port}

	// Start the HTTP server in a separate goroutine
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Listen for interrupt or termination signals to gracefully stop the server
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received
	<-signalChan

	return s.Stop()
}

func (s *Server) Stop() error {
	log.Println("Shutting down server gracefully...")

	// Shut down the HTTP server gracefully
	if err := s.httpServer.Shutdown(nil); err != nil {
		return err
	}

	log.Println("Server shut down gracefully.")
	return nil
}
