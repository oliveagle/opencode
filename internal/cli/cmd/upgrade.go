package cmd

import (
	"fmt"

	"github.com/anomalyco/opencode/internal/installation"
	"github.com/spf13/cobra"
)

func NewUpgradeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade OpenCode to the latest version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Current version: %s\n", installation.Version)
			// TODO: Implement upgrade logic
			fmt.Println("Upgrade functionality will be implemented")
		},
	}
	return cmd
}
