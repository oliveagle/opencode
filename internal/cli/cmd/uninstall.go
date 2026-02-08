package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewUninstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall OpenCode",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("Are you sure you want to uninstall OpenCode? [y/N]: ")
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			if response == "y\n" || response == "Y\n" {
				// TODO: Implement uninstall logic
				fmt.Println("Uninstall functionality will be implemented")
			} else {
				fmt.Println("Uninstall cancelled")
			}
		},
	}
	return cmd
}
