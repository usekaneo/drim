package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/usekaneo/drim/pkg/docker"
	"github.com/usekaneo/drim/pkg/ui"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Edit Kaneo configuration",
	Long: `Opens .env in your default editor ($EDITOR or nano). 
After saving, services are restarted automatically.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		envPath := ".env"

		if _, err := os.Stat(envPath); os.IsNotExist(err) {
			return fmt.Errorf(".env file not found. Run 'drim setup' first")
		}

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "nano"
		}

		ui.Info(fmt.Sprintf("Opening %s in %s...", envPath, editor))

		editorCmd := exec.Command(editor, envPath)
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr

		if err := editorCmd.Run(); err != nil {
			return fmt.Errorf("failed to open editor: %w", err)
		}

		shouldRestart, err := ui.Confirm("Configuration updated. Restart services to apply changes?")
		if err != nil {
			return err
		}

		if shouldRestart {
			ui.Info("Restarting services...")
			if err := docker.ComposeRestart(); err != nil {
				return fmt.Errorf("failed to restart services: %w", err)
			}
			ui.Success("✨ Services restarted successfully!")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
}
