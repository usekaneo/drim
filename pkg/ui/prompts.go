package ui

import (
	"github.com/usekaneo/drim/pkg/generator"
)

func PromptSetupConfig() (*generator.Config, error) {
	config := generator.NewDefaultConfig()

	Info("\n📝 Configuration")
	Info("Press Enter to use default values\n")

	useCaddy, err := Confirm("Do you want to use a reverse proxy? (recommended)")
	if err != nil {
		return nil, err
	}
	config.UseCaddy = useCaddy

	if !config.UseCaddy {
		Warning("You chose not to use a reverse proxy.")
		Info("API will be exposed on port 1337, Web on port 5173")
	}

	domain, err := Prompt("Enter your domain (e.g., kaneo.example.com) [optional]:")
	if err != nil {
		return nil, err
	}
	config.Domain = domain

	if config.Domain == "" {
		if config.UseCaddy {
			Info("No domain specified. Kaneo will be accessible at http://localhost")
		} else {
			Info("No domain specified. Services will be available at:")
			Info("  - API: http://localhost:1337")
			Info("  - Web: http://localhost:5173")
		}
	} else {
		if config.UseCaddy {
			Info("Kaneo will be accessible at https://" + config.Domain)
			Info("Make sure your domain's DNS A record points to this server's IP address")
		} else {
			Info("Services will be available. Configure your reverse proxy for: " + config.Domain)
		}
	}

	return config, nil
}
