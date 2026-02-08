package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/anomalyco/opencode/pkg/tui"
)

func main() {
	var serverURL = flag.String("server", "http://localhost:8080", "Backend server URL")
	flag.Parse()

	app, err := tui.NewApp(*serverURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}

