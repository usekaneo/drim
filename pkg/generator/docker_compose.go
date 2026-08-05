package generator

import (
	"os"
)

func GenerateDockerCompose(config *Config) error {
	var content string

	if config.UseCaddy {
		content = `services:
  postgres:
    image: postgres:16-alpine
    env_file: .env
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped
    networks:
      - kaneo
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U kaneo -d kaneo"]
      interval: 10s
      timeout: 5s
      retries: 5

  kaneo:
    image: ghcr.io/usekaneo/kaneo:latest
    env_file: .env
    depends_on:
      postgres:
        condition: service_healthy
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://127.0.0.1:5173/api/health"]
      interval: 30s
      timeout: 5s
      start_period: 60s
      retries: 3
    networks:
      - kaneo

  caddy:
    image: caddy:2-alpine
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile:ro
      - caddy_data:/data
      - caddy_config:/config
    depends_on:
      - kaneo
    networks:
      - kaneo

volumes:
  postgres_data:
  caddy_data:
  caddy_config:

networks:
  kaneo:
    driver: bridge
`
	} else {
		content = `services:
  postgres:
    image: postgres:16-alpine
    env_file: .env
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped
    networks:
      - kaneo
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U kaneo -d kaneo"]
      interval: 10s
      timeout: 5s
      retries: 5

  kaneo:
    image: ghcr.io/usekaneo/kaneo:latest
    env_file: .env
    depends_on:
      postgres:
        condition: service_healthy
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://127.0.0.1:5173/api/health"]
      interval: 30s
      timeout: 5s
      start_period: 60s
      retries: 3
    ports:
      - "5173:5173"
    networks:
      - kaneo

volumes:
  postgres_data:

networks:
  kaneo:
    driver: bridge
`
	}

	return os.WriteFile("docker-compose.yml", []byte(content), 0644)
}
