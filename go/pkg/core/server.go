package core

import (
	"fmt"
	"log"
)

type Server struct {
	// TODO: add server fields
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start() error {
	log.Println("Starting opencode backend server")
	fmt.Println("Backend server starting on :8080")
	// TODO: implement server start logic
	return nil
}

func (s *Server) Stop() error {
	log.Println("Stopping opencode backend server")
	return nil
}
