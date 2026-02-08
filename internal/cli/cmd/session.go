package cmd

import (
	"github.com/spf13/cobra"
)

func NewSessionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "session",
		Short: "Manage sessions",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List sessions",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement session list
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "delete",
		Short: "Delete a session",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement session delete
		},
	})
	return cmd
}
