package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/usekaneo/drim/pkg/banner"
	"github.com/usekaneo/drim/pkg/docker"
	"github.com/usekaneo/drim/pkg/generator"
	"github.com/usekaneo/drim/pkg/ui"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Deploy the entire Kaneo stack",
	Long: `Deploys the entire stack including PostgreSQL, Kaneo API, Kaneo Web, and Caddy.
	
- Installs Docker if missing (on supported systems)
- Generates docker-compose.yml, Caddyfile, and .env
- Pulls images and starts everything with docker compose up -d
- Configures Caddy automatically (HTTPS on 443, HTTP on 80)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		banner.Print()

		ui.Info("🚀 Starting Kaneo setup...")

		ui.Info("Checking Docker installation...")
		if !docker.IsInstalled() {
			ui.Warning("Docker is not installed.")
			ui.Info("Attempting to install Docker...")
			if err := docker.Install(); err != nil {
				return fmt.Errorf("failed to install Docker: %w", err)
			}
			ui.Success("Docker installed successfully!")
		} else {
			ui.Success("Docker is already installed")
		}

		if !docker.IsComposeAvailable() {
			return fmt.Errorf("Docker Compose is not available")
		}
		ui.Success("Docker Compose is available")

		config, err := ui.PromptSetupConfig()
		if err != nil {
			return fmt.Errorf("failed to get configuration: %w", err)
		}

		ui.Info("Generating configuration files...")

		if err := generator.GenerateDockerCompose(config); err != nil {
			return fmt.Errorf("failed to generate docker-compose.yml: %w", err)
		}
		ui.Success("Generated docker-compose.yml")

		if config.UseCaddy {
			if err := generator.GenerateCaddyfile(config); err != nil {
				return fmt.Errorf("failed to generate Caddyfile: %w", err)
			}
			ui.Success("Generated Caddyfile")
		}

		if err := generator.GenerateEnvFile(config); err != nil {
			return fmt.Errorf("failed to generate .env: %w", err)
		}
		ui.Success("Generated .env")

		ui.Info("Pulling Docker images (this may take a few minutes)...")
		if err := docker.ComposePull(); err != nil {
			return fmt.Errorf("failed to pull images: %w", err)
		}
		ui.Success("Images pulled successfully")

		ui.Info("Starting services...")
		if err := docker.ComposeUp(); err != nil {
			return fmt.Errorf("failed to start services: %w", err)
		}

		ui.Success("\n✨ Kaneo is now running!")
		if config.UseCaddy {
			if config.Domain != "" {
				ui.Info(fmt.Sprintf("🌐 Access your instance at: https://%s", config.Domain))
				ui.Info("(HTTPS certificate will be generated automatically)")
			} else {
				ui.Info("🌐 Access your instance at: http://localhost")
			}
		} else {
			ui.Info("🌐 Services are running:")
			ui.Info("   • API: http://localhost:1337")
			ui.Info("   • Web: http://localhost:5173")
			if config.Domain != "" {
				ui.Info(fmt.Sprintf("\n📝 Configure your reverse proxy to forward %s to these services", config.Domain))
			} else {
				ui.Info("\n📝 Set up your reverse proxy to forward requests to these services")
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
