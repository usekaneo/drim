package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/usekaneo/drim/pkg/banner"
	"github.com/usekaneo/drim/pkg/docker"
	"github.com/usekaneo/drim/pkg/generator"
	"github.com/usekaneo/drim/pkg/ui"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Deploy the entire Kaneo stack",
	Long: `Deploys the entire stack including PostgreSQL, Kaneo (unified web and API), and Caddy.
	
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

			added, err := docker.AddCurrentUserToDockerGroup()
			if err != nil {
				return fmt.Errorf("failed to configure Docker access: %w", err)
			}
			// Root already reaches the socket, so only an unprivileged session
			// has to re-login before the new group membership applies.
			if added && os.Geteuid() != 0 {
				ui.Warning("Docker access was enabled for your user.")
				ui.Info("Log out and back in, then rerun 'drim setup'.")
				return nil
			}
		} else {
			ui.Success("Docker is already installed")
		}

		if !docker.IsComposeAvailable() {
			return fmt.Errorf("Docker Compose is not available")
		}
		ui.Success("Docker Compose is available")

		if !docker.IsDaemonReachable() {
			ui.Error("Cannot talk to the Docker daemon.")
			if os.Geteuid() != 0 {
				ui.Info("Add yourself to the docker group, then log out and back in:")
				ui.Info("  sudo usermod -aG docker $USER")
				ui.Info("Or rerun this command with sudo.")
			} else {
				ui.Info("Check that the Docker daemon is running: systemctl status docker")
			}
			return fmt.Errorf("docker daemon is not reachable")
		}

		config, err := setupConfig(cmd)
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
			ui.Info("🌐 Kaneo is running at http://localhost:5173")
			if config.Domain != "" {
				ui.Info(fmt.Sprintf("\n📝 Configure your reverse proxy to forward %s to Kaneo on port 5173", config.Domain))
			} else {
				ui.Info("\n📝 Set up your reverse proxy to forward requests to Kaneo on port 5173")
			}
		}

		return nil
	},
}

// setupConfig skips the prompts when the caller passed configuration as
// flags, so setup works from a non-interactive shell such as install.sh.
func setupConfig(cmd *cobra.Command) (*generator.Config, error) {
	domain, err := cmd.Flags().GetString("domain")
	if err != nil {
		return nil, err
	}
	noReverseProxy, err := cmd.Flags().GetBool("no-reverse-proxy")
	if err != nil {
		return nil, err
	}

	if !cmd.Flags().Changed("domain") && !cmd.Flags().Changed("no-reverse-proxy") {
		return ui.PromptSetupConfig()
	}

	config := generator.NewDefaultConfig()
	config.Domain = domain
	config.UseCaddy = !noReverseProxy
	ui.Info(fmt.Sprintf("Using configuration from flags (domain: %q, reverse proxy: %t)", config.Domain, config.UseCaddy))
	return config, nil
}

func init() {
	setupCmd.Flags().String("domain", "", "Domain to serve Kaneo on, for example kaneo.example.com")
	setupCmd.Flags().Bool("no-reverse-proxy", false, "Skip Caddy and expose Kaneo directly on port 5173")
	rootCmd.AddCommand(setupCmd)
}
