package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/usekaneo/drim/pkg/docker"
	"github.com/usekaneo/drim/pkg/ui"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart all Kaneo services",
	Long:  `Restarts all Kaneo services using Docker Compose.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.Info("Restarting Kaneo services...")
		
		if err := docker.ComposeRestart(); err != nil {
			return fmt.Errorf("failed to restart services: %w", err)
		}
		
		ui.Success("✨ Kaneo services restarted successfully!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}


