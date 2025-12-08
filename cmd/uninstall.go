package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/usekaneo/drim/pkg/docker"
	"github.com/usekaneo/drim/pkg/ui"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall Kaneo",
	Long: `Stops and removes containers, networks, and volumes. 
Optionally removes images.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.Warning("⚠️  This will remove all Kaneo containers, networks, and volumes.")
		ui.Warning("Your data will be permanently deleted!")
		
		confirmed, err := ui.Confirm("Are you sure you want to uninstall Kaneo?")
		if err != nil {
			return err
		}
		
		if !confirmed {
			ui.Info("Uninstall cancelled.")
			return nil
		}
		
		ui.Info("Stopping and removing containers...")
		if err := docker.ComposeDown(true); err != nil {
			return fmt.Errorf("failed to remove containers: %w", err)
		}
		ui.Success("Containers and volumes removed")
		
		removeImages, err := ui.Confirm("Do you also want to remove downloaded images?")
		if err != nil {
			return err
		}
		
		if removeImages {
			ui.Info("Removing images...")
			if err := docker.RemoveImages(); err != nil {
				ui.Warning(fmt.Sprintf("Failed to remove some images: %v", err))
			} else {
				ui.Success("Images removed")
			}
		}
		
		ui.Success("✨ Kaneo uninstalled successfully!")
		ui.Info("Configuration files (docker-compose.yml, Caddyfile, .env) have been kept.")
		ui.Info("You can delete them manually if needed.")
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}


