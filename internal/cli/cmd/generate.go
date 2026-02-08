package cmd

import (
	"github.com/spf13/cobra"
)

func NewGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate code using AI",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement generate command
		},
	}
	return cmd
}
