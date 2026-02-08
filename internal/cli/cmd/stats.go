package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewStatsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Show usage statistics",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement stats command
			fmt.Println("Usage statistics will be shown here")
		},
	}
	return cmd
}
