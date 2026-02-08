package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var servePort int
var serveHost string

func NewServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the OpenCode server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Starting OpenCode server on %s:%d\n", serveHost, servePort)
			// TODO: Implement server logic
		},
	}
	cmd.Flags().IntVarP(&servePort, "port", "p", 3000, "Port to listen on")
	cmd.Flags().StringVar(&serveHost, "host", "localhost", "Host to bind to")
	return cmd
}
