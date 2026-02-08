package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewModelsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "models",
		Short: "List and manage AI models",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List available models",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement models list
			fmt.Println("Available models will be listed here")
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "pull",
		Short: "Pull model information from models.dev",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement models pull
			fmt.Println("Pulling models from models.dev")
		},
	})
	return cmd
}
