package cmd

import (
	"github.com/spf13/cobra"
)

func NewDebugCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "debug",
		Short: "Debug commands",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "config",
		Short: "Show current configuration",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement debug config
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "lsp",
		Short: "Debug LSP",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement debug lsp
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "file",
		Short: "Debug file operations",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement debug file
		},
	})
	return cmd
}
