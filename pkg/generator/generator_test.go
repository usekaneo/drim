package generator

import (
	"os"
	"testing"
)

func TestNewDefaultConfig(t *testing.T) {
	config := NewDefaultConfig()

	if config.PostgresUser != "kaneo" {
		t.Errorf("Expected PostgresUser to be 'kaneo', got '%s'", config.PostgresUser)
	}

	if config.PostgresDB != "kaneo" {
		t.Errorf("Expected PostgresDB to be 'kaneo', got '%s'", config.PostgresDB)
	}

	if config.PostgresPass == "" {
		t.Error("Expected PostgresPass to be generated")
	}

	if config.AuthSecret == "" {
		t.Error("Expected AuthSecret to be generated")
	}

	if len(config.PostgresPass) != 32 {
		t.Errorf("Expected PostgresPass length to be 32, got %d", len(config.PostgresPass))
	}

	if len(config.AuthSecret) != 64 {
		t.Errorf("Expected AuthSecret length to be 64, got %d", len(config.AuthSecret))
	}
}

func TestGenerateDockerCompose(t *testing.T) {
	config := NewDefaultConfig()

	// Create a temporary file
	tmpFile := "test-docker-compose.yml"
	defer os.Remove(tmpFile)

	// Override the file path for testing
	originalFile := "docker-compose.yml"
	defer func() {
		os.Remove(originalFile)
	}()

	err := GenerateDockerCompose(config)
	if err != nil {
		t.Fatalf("Failed to generate docker-compose.yml: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat("docker-compose.yml"); os.IsNotExist(err) {
		t.Error("docker-compose.yml was not created")
	}

	// Read and verify content
	content, err := os.ReadFile("docker-compose.yml")
	if err != nil {
		t.Fatalf("Failed to read docker-compose.yml: %v", err)
	}

	contentStr := string(content)

	// Check for required services
	requiredStrings := []string{
		"postgres:",
		"kaneo:",
		"caddy:",
		"postgres:16-alpine",
		"ghcr.io/usekaneo/kaneo:latest",
		"caddy:2-alpine",
	}

	for _, required := range requiredStrings {
		if !contains(contentStr, required) {
			t.Errorf("docker-compose.yml missing required string: %s", required)
		}
	}

	for _, legacyImage := range []string{
		"ghcr.io/usekaneo/api:latest",
		"ghcr.io/usekaneo/web:latest",
	} {
		if contains(contentStr, legacyImage) {
			t.Errorf("docker-compose.yml still uses legacy image: %s", legacyImage)
		}
	}
}

func TestGenerateDockerComposeWithoutCaddy(t *testing.T) {
	config := NewDefaultConfig()
	config.UseCaddy = false
	defer os.Remove("docker-compose.yml")

	if err := GenerateDockerCompose(config); err != nil {
		t.Fatalf("Failed to generate docker-compose.yml: %v", err)
	}

	content, err := os.ReadFile("docker-compose.yml")
	if err != nil {
		t.Fatalf("Failed to read docker-compose.yml: %v", err)
	}
	contentStr := string(content)

	for _, expected := range []string{
		"kaneo:",
		"ghcr.io/usekaneo/kaneo:latest",
		`"5173:5173"`,
	} {
		if !contains(contentStr, expected) {
			t.Errorf("docker-compose.yml missing required string: %s", expected)
		}
	}

	for _, unexpected := range []string{
		"caddy:",
		"ghcr.io/usekaneo/api:latest",
		"ghcr.io/usekaneo/web:latest",
		`"1337:1337"`,
	} {
		if contains(contentStr, unexpected) {
			t.Errorf("docker-compose.yml contains unexpected string: %s", unexpected)
		}
	}
}

func TestGenerateCaddyfile(t *testing.T) {
	tests := []struct {
		name     string
		domain   string
		expected []string
	}{
		{
			name:   "With domain",
			domain: "kaneo.example.com",
			expected: []string{
				"kaneo.example.com",
				"reverse_proxy http://kaneo:5173",
			},
		},
		{
			name:   "Without domain",
			domain: "",
			expected: []string{
				"auto_https off",
				":80",
				"reverse_proxy http://kaneo:5173",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewDefaultConfig()
			config.Domain = tt.domain

			defer os.Remove("Caddyfile")

			err := GenerateCaddyfile(config)
			if err != nil {
				t.Fatalf("Failed to generate Caddyfile: %v", err)
			}

			// Check if file exists
			if _, err := os.Stat("Caddyfile"); os.IsNotExist(err) {
				t.Error("Caddyfile was not created")
			}

			// Read and verify content
			content, err := os.ReadFile("Caddyfile")
			if err != nil {
				t.Fatalf("Failed to read Caddyfile: %v", err)
			}

			contentStr := string(content)

			for _, expected := range tt.expected {
				if !contains(contentStr, expected) {
					t.Errorf("Caddyfile missing expected string: %s", expected)
				}
			}

			for _, legacyTarget := range []string{"http://api:1337", "http://web:5173"} {
				if contains(contentStr, legacyTarget) {
					t.Errorf("Caddyfile still targets legacy service: %s", legacyTarget)
				}
			}
		})
	}
}

func TestGenerateEnvFile(t *testing.T) {
	config := NewDefaultConfig()
	config.Domain = "kaneo.example.com"

	defer os.Remove(".env")

	err := GenerateEnvFile(config)
	if err != nil {
		t.Fatalf("Failed to generate .env: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		t.Error(".env was not created")
	}

	// Read and verify content
	content, err := os.ReadFile(".env")
	if err != nil {
		t.Fatalf("Failed to read .env: %v", err)
	}

	contentStr := string(content)

	// Check for required variables
	requiredStrings := []string{
		"POSTGRES_USER=",
		"POSTGRES_PASSWORD=",
		"POSTGRES_DB=",
		"DATABASE_URL=",
		"AUTH_SECRET=",
		"KANEO_CLIENT_URL=",
		"KANEO_API_URL=",
		"DOMAIN=",
	}

	for _, required := range requiredStrings {
		if !contains(contentStr, required) {
			t.Errorf(".env missing required variable: %s", required)
		}
	}

	// Verify domain is included
	if !contains(contentStr, "DOMAIN=kaneo.example.com") {
		t.Error(".env does not contain the correct domain")
	}
}

func TestGenerateEnvFileWithoutCaddyUsesUnifiedOrigin(t *testing.T) {
	config := NewDefaultConfig()
	config.UseCaddy = false
	defer os.Remove(".env")

	if err := GenerateEnvFile(config); err != nil {
		t.Fatalf("Failed to generate .env: %v", err)
	}

	content, err := os.ReadFile(".env")
	if err != nil {
		t.Fatalf("Failed to read .env: %v", err)
	}
	contentStr := string(content)

	for _, expected := range []string{
		"KANEO_CLIENT_URL=http://localhost:5173",
		"KANEO_API_URL=http://localhost:5173/api",
	} {
		if !contains(contentStr, expected) {
			t.Errorf(".env missing unified origin: %s", expected)
		}
	}
	if contains(contentStr, "localhost:1337") {
		t.Error(".env still exposes the legacy API origin")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && (s[0:len(substr)] == substr || contains(s[1:], substr))))
}
