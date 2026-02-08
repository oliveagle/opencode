package cmd

import (
	"github.com/spf13/cobra"
)

func NewRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run [prompt]",
		Short: "Run opencode with a prompt",
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement run command
			// This will start the TUI or run in headless mode
		},
	}
	return cmd
}
