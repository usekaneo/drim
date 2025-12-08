package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/usekaneo/drim/pkg/docker"
	"github.com/usekaneo/drim/pkg/ui"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop all Kaneo services",
	Long:  `Stops all Kaneo services using Docker Compose.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.Info("Stopping Kaneo services...")
		
		if err := docker.ComposeStop(); err != nil {
			return fmt.Errorf("failed to stop services: %w", err)
		}
		
		ui.Success("✨ Kaneo services stopped successfully!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}



