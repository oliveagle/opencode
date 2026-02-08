package cli

import (
	"fmt"
	"os"

	"github.com/anomalyco/opencode/internal/cli/cmd"
	"github.com/anomalyco/opencode/internal/installation"
	"github.com/anomalyco/opencode/internal/util/log"
	"github.com/spf13/cobra"
)

var (
	printLogs bool
	logLevel  string
)

var rootCmd = &cobra.Command{
	Use:     "opencode",
	Short:   "The open source AI coding agent",
	Version: installation.Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize logging
		level := log.LevelInfo
		if logLevel != "" {
			switch logLevel {
			case "DEBUG":
				level = log.LevelDebug
			case "INFO":
				level = log.LevelInfo
			case "WARN":
				level = log.LevelWarn
			case "ERROR":
				level = log.LevelError
			}
		} else if installation.IsLocal() {
			level = log.LevelDebug
		}

		log.Init(log.Config{
			Print: printLogs,
			Dev:   installation.IsLocal(),
			Level: level,
		})

		// Set environment flags
		os.Setenv("AGENT", "1")
		os.Setenv("OPENCODE", "1")

		log.Default.Info("opencode", map[string]interface{}{
			"version": installation.Version,
			"args":    os.Args[1:],
		})
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&printLogs, "print-logs", false, "print logs to stderr")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "", "log level (DEBUG, INFO, WARN, ERROR)")

	// Add subcommands
	rootCmd.AddCommand(cmd.NewRunCommand())
	rootCmd.AddCommand(cmd.NewGenerateCommand())
	rootCmd.AddCommand(cmd.NewAuthCommand())
	rootCmd.AddCommand(cmd.NewAgentCommand())
	rootCmd.AddCommand(cmd.NewUpgradeCommand())
	rootCmd.AddCommand(cmd.NewUninstallCommand())
	rootCmd.AddCommand(cmd.NewServeCommand())
	rootCmd.AddCommand(cmd.NewModelsCommand())
	rootCmd.AddCommand(cmd.NewStatsCommand())
	rootCmd.AddCommand(cmd.NewMcpCommand())
	rootCmd.AddCommand(cmd.NewExportCommand())
	rootCmd.AddCommand(cmd.NewImportCommand())
	rootCmd.AddCommand(cmd.NewSessionCommand())
	rootCmd.AddCommand(cmd.NewDebugCommand())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
