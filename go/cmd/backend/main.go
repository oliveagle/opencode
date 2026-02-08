package main

import (
	"fmt"
	"os"

	"github.com/anomalyco/opencode/pkg/core"
)

func main() {
	server := core.NewServer()
	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting server: %v\n", err)
		os.Exit(1)
	}
}
