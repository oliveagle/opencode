package cmd

import (
	"github.com/spf13/cobra"
)

func NewImportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import session data",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement import command
		},
	}
	return cmd
}
