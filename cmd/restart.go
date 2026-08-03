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

		force, err := cmd.Flags().GetBool("force")
		if err == nil && force {
			ui.Info("Force restarting services...")
			if err := docker.ComposeUpD(); err != nil {
				return fmt.Errorf("failed to force restart services: %w", err)
			}
			ui.Success("✨ Services force restarted successfully!")
			return nil
		}

		if err := docker.ComposeRestart(); err != nil {
			return fmt.Errorf("failed to restart services: %w", err)
		}

		ui.Success("✨ Kaneo services restarted successfully!")
		return nil
	},
}

func init() {
	restartCmd.Flags().Bool("force", false, "Run compose up -d")
	rootCmd.AddCommand(restartCmd)
}
