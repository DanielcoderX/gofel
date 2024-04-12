package rpc

import (
	"log"
	"net/http"

	"github.com/DanielcoderX/gofel"
)

type Server struct {
	config *config.Config
}

func NewServer(cfg *config.Config) (*Server, error) {
	return &Server{
		config: cfg,
	}, nil
}

func (s *Server) Start() error {
	http.HandleFunc(s.config.Path, handleWebSocket)
	log.Printf("RPC starting on %s...", s.config.Port)
	return http.ListenAndServe(":"+s.config.Port, nil)
}
