package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/usekaneo/drim/pkg/docker"
	"github.com/usekaneo/drim/pkg/ui"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start all Kaneo services",
	Long:  `Starts all Kaneo services using Docker Compose.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireComposeContext(); err != nil {
			return err
		}

		ui.Info("Starting Kaneo services...")

		if err := docker.ComposeUp(); err != nil {
			return fmt.Errorf("failed to start services: %w", err)
		}

		ui.Success("✨ Kaneo services started successfully!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
