package cmd

import (
	"github.com/spf13/cobra"
)

func NewAuthCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "login",
		Short: "Login to OpenCode",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement auth login
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "logout",
		Short: "Logout from OpenCode",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement auth logout
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Show authentication status",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement auth status
		},
	})
	return cmd
}
