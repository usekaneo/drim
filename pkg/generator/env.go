package generator

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
)

func GenerateEnvFile(config *Config) error {
	clientURL := "http://localhost"
	apiURL := "http://localhost/api"

	if config.Domain != "" {
		clientURL = fmt.Sprintf("https://%s", config.Domain)
		apiURL = fmt.Sprintf("https://%s/api", config.Domain)
	}

	content := fmt.Sprintf(`KANEO_CLIENT_URL=%s
KANEO_API_URL=%s

DATABASE_URL=postgresql://%s:%s@postgres:5432/%s
POSTGRES_DB=%s
POSTGRES_USER=%s
POSTGRES_PASSWORD=%s

AUTH_SECRET=%s

DOMAIN=%s
`,
		clientURL,
		apiURL,
		config.PostgresUser,
		config.PostgresPass,
		config.PostgresDB,
		config.PostgresDB,
		config.PostgresUser,
		config.PostgresPass,
		config.AuthSecret,
		config.Domain,
	)

	return os.WriteFile(".env", []byte(content), 0600)
}

func generateRandomPassword(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "changeme" + fmt.Sprintf("%d", length)
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}
