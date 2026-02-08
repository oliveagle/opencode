package cmd

import (
	"github.com/spf13/cobra"
)

func NewExportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export session data",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement export command
		},
	}
	return cmd
}
