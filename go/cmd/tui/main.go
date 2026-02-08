package main

import (
	"fmt"
	"os"

	"github.com/anomalyco/opencode/pkg/tui"
)

func main() {
	app := tui.NewApp()
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
