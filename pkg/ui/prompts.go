package ui

import (
	"github.com/usekaneo/drim/pkg/generator"
)

func PromptSetupConfig() (*generator.Config, error) {
	config := generator.NewDefaultConfig()

	Info("\n📝 Configuration")
	Info("Press Enter to skip optional fields\n")

	domain, err := Prompt("Enter your domain (e.g., kaneo.example.com) [optional]:")
	if err != nil {
		return nil, err
	}
	config.Domain = domain

	if config.Domain == "" {
		Info("No domain specified. Kaneo will be accessible at http://localhost")
	} else {
		Info("Kaneo will be accessible at https://" + config.Domain)
		Info("Make sure your domain's DNS A record points to this server's IP address")
	}

	return config, nil
}
