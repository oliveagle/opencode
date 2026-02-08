package cmd

import (
	"github.com/spf13/cobra"
)

func NewAgentCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Manage agents",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List available agents",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement agent list
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "create",
		Short: "Create a new agent",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement agent create
		},
	})
	return cmd
}
