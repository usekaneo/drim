package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/usekaneo/drim/pkg/docker"
	"github.com/usekaneo/drim/pkg/ui"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade Kaneo to the latest version",
	Long: `Pulls the latest Kaneo images and restarts only the changed services.
	
This command will:
- Pull the latest images from the registry
- Restart services that have been updated`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.Info("Pulling latest Kaneo images...")
		
		if err := docker.ComposePull(); err != nil {
			return fmt.Errorf("failed to pull images: %w", err)
		}
		ui.Success("Images pulled successfully")
		
		ui.Info("Restarting services with new images...")
		if err := docker.ComposeUp(); err != nil {
			return fmt.Errorf("failed to restart services: %w", err)
		}
		
		ui.Success("✨ Kaneo upgraded successfully!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}




