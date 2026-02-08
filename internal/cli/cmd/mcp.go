package cmd

import (
	"github.com/spf13/cobra"
)

func NewMcpCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Manage MCP (Model Context Protocol) servers",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List configured MCP servers",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement MCP list
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "add",
		Short: "Add a new MCP server",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement MCP add
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "remove",
		Short: "Remove an MCP server",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement MCP remove
		},
	})
	return cmd
}
