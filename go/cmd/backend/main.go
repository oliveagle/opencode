package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/anomalyco/opencode/pkg/api"
	"github.com/anomalyco/opencode/pkg/core"
)

func main() {
	var (
		port    = flag.String("port", "8080", "Server port")
		baseDir = flag.String("dir", ".", "Base directory for file operations")
	)
	flag.Parse()

	// Initialize services
	fileSvc := core.NewFileService(*baseDir)
	projSvc, err := core.NewProjectService(os.ExpandEnv("$HOME/.opencode"))
	if err != nil {
		log.Fatalf("Failed to initialize project service: %v", err)
	}

	// Setup HTTP handlers
	mux := http.NewServeMux()
	handler := api.NewHandler(fileSvc, projSvc)
	handler.RegisterRoutes(mux)

	// Start server
	server := &http.Server{
		Addr:    ":" + *port,
		Handler: mux,
	}

	go func() {
		log.Printf("Starting server on http://localhost:%s", *port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
	server.Close()
}

