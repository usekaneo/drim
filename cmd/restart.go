package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/usekaneo/drim/pkg/docker"
	"github.com/usekaneo/drim/pkg/ui"
)

var recreate bool

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart all Kaneo services",
	Long: `Restarts all Kaneo services using Docker Compose.

Use --recreate to recreate the containers instead, which is required for
changes to .env to take effect.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireComposeContext(); err != nil {
			return err
		}

		ui.Info("Restarting Kaneo services...")

		successMessage := "✨ Kaneo services restarted successfully!"
		if recreate {
			ui.Info("Recreating services to apply configuration changes...")
			if err := docker.ComposeUp(); err != nil {
				return fmt.Errorf("failed to recreate services: %w", err)
			}
			successMessage = "✨ Kaneo services recreated successfully!"
		} else {
			if err := docker.ComposeRestart(); err != nil {
				return fmt.Errorf("failed to restart services: %w", err)
			}
		}

		ui.Success(successMessage)
		return nil
	},
}

func init() {
	restartCmd.Flags().BoolVar(&recreate, "recreate", false, "Recreate containers so .env changes take effect")
	rootCmd.AddCommand(restartCmd)
}
