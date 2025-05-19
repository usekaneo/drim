package cmd

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	domain          string
	jwtToken        string
	proxyType       string
	useHttps        bool
	disableRegister bool
)

var (
	success = color.New(color.FgGreen).SprintFunc()
	info    = color.New(color.FgCyan).SprintFunc()
	warning = color.New(color.FgYellow).SprintFunc()
	errorC  = color.New(color.FgRed).SprintFunc()
	bold    = color.New(color.Bold).SprintFunc()
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy Kaneo using Docker",
	Long: `Deploy Kaneo using Docker and Docker Compose.
This command will create a Docker Compose file and start the containers.
You can specify your domain and other configuration options.`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		// Interactive mode if no domain is provided
		if domain == "" && (proxyType == "traefik" || proxyType == "nginx") {
			fmt.Printf("%s Please enter your domain name: ", info("?"))
			userDomain, _ := reader.ReadString('\n')
			domain = strings.TrimSpace(userDomain)
			if domain == "" {
				fmt.Printf("%s Domain is required for %s proxy. Using 'none' as default.\n", warning("⚠"), proxyType)
				proxyType = "none"
			}
		}

		// If proxy type is not specified or invalid, ask the user
		if proxyType != "traefik" && proxyType != "nginx" && proxyType != "none" {
			fmt.Printf("%s Invalid proxy type. Please choose one of the following:\n", warning("⚠"))
			fmt.Printf("  1. %s - For production with separate domains (api.domain.com and domain.com)\n", bold("traefik"))
			fmt.Printf("  2. %s - For production with path-based routing (domain.com/api)\n", bold("nginx"))
			fmt.Printf("  3. %s - For local development (localhost:5173 and localhost:1337)\n", bold("none"))
			fmt.Printf("%s Enter your choice (1-3): ", info("?"))

			choice, _ := reader.ReadString('\n')
			choice = strings.TrimSpace(choice)

			switch choice {
			case "1":
				proxyType = "traefik"
				if domain == "" {
					fmt.Printf("%s Please enter your domain name: ", info("?"))
					userDomain, _ := reader.ReadString('\n')
					domain = strings.TrimSpace(userDomain)
					if domain == "" {
						fmt.Printf("%s Domain is required for Traefik. Using 'none' as default.\n", warning("⚠"))
						proxyType = "none"
					}
				}
			case "2":
				proxyType = "nginx"
				if domain == "" {
					fmt.Printf("%s Please enter your domain name: ", info("?"))
					userDomain, _ := reader.ReadString('\n')
					domain = strings.TrimSpace(userDomain)
					if domain == "" {
						fmt.Printf("%s Domain is required for Nginx. Using 'none' as default.\n", warning("⚠"))
						proxyType = "none"
					}
				}
			default:
				proxyType = "none"
				fmt.Printf("%s Using 'none' as default proxy type.\n", info("ℹ"))
			}
		}

		if useHttps && proxyType == "none" {
			fmt.Printf("%s HTTPS is only available with Traefik or Nginx. Disabling HTTPS.\n", warning("⚠"))
			useHttps = false
		}

		if jwtToken == "" {
			jwtToken = generateRandomToken(32)
			fmt.Printf("%s Generated JWT token: %s\n", info("ℹ"), jwtToken)
		}

		switch proxyType {
		case "traefik":
			createTraefikComposeFile()
		case "nginx":
			createNginxComposeFile()
		case "none":
			createBasicComposeFile()
		default:
			fmt.Printf("%s Invalid proxy type. Using 'none' as default.\n", warning("⚠"))
			createBasicComposeFile()
		}

		startContainers()
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringVarP(&domain, "domain", "d", "", "Your domain name (e.g., kaneo.example.com)")
	deployCmd.Flags().StringVarP(&jwtToken, "jwt", "j", "", "JWT access token (will be generated if not provided)")
	deployCmd.Flags().StringVarP(&proxyType, "proxy", "p", "", "Proxy type: traefik, nginx, or none")
	deployCmd.Flags().BoolVarP(&useHttps, "https", "s", false, "Enable HTTPS (only applicable with traefik or nginx)")
	deployCmd.Flags().BoolVarP(&disableRegister, "disable-register", "r", false, "Disable user registration")
}

func generateRandomToken(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Printf("%s Error generating random token: %v\n", errorC("✗"), err)
		os.Exit(1)
	}
	return base64.URLEncoding.EncodeToString(b)[:length]
}

func createBasicComposeFile() {
	content := `services:
  backend:
    image: ghcr.io/usekaneo/api:latest
    environment:
      JWT_ACCESS: "` + jwtToken + `"
      DB_PATH: "/app/apps/api/data/kaneo.db"
      DISABLE_REGISTRATION: "` + fmt.Sprintf("%t", disableRegister) + `"
    ports:
      - 1337:1337
    restart: unless-stopped
    volumes:
      - sqlite_data:/app/apps/api/data

  frontend:
    image: ghcr.io/usekaneo/web:latest
    environment:
      KANEO_API_URL: "http://localhost:1337"
    ports:
      - 5173:5173
    restart: unless-stopped

volumes:
  sqlite_data:
`
	writeToFile("compose.yml", content)
	fmt.Printf("%s Created basic Docker Compose file: %s\n", success("✓"), bold("compose.yml"))
	fmt.Printf("%s After deployment, Kaneo will be available at: %s\n", info("ℹ"), bold("http://localhost:5173"))
}

func createTraefikComposeFile() {
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")

	apiDomain := "api-" + domain

	protocol := "http"
	if useHttps {
		protocol = "https"
	}

	disableRegistrationValue := "false"
	if disableRegister {
		disableRegistrationValue = "true"
	}

	content := `services:
  traefik:
    image: "traefik:v3.3"
    container_name: "traefik"
    command:
      - "--providers.docker=true"
      - "--entryPoints.web.address=:80"
`
	if useHttps {
		content += `      - "--entryPoints.websecure.address=:443"
      - "--certificatesresolvers.myresolver.acme.httpchallenge=true"
      - "--certificatesresolvers.myresolver.acme.httpchallenge.entrypoint=web"
      - "--certificatesresolvers.myresolver.acme.email=admin@` + domain + `"
      - "--certificatesresolvers.myresolver.acme.storage=/letsencrypt/acme.json"
`
	}

	content += `    ports:
      - "80:80"
`
	if useHttps {
		content += `      - "443:443"
`
	}

	content += `    volumes:
      - "/var/run/docker.sock:/var/run/docker.sock:ro"
`
	if useHttps {
		content += `      - "letsencrypt:/letsencrypt"
`
	}

	content += `    networks:
      - traefik-net

  backend:
    image: ghcr.io/usekaneo/api:latest
    environment:
      JWT_ACCESS: "` + jwtToken + `"
      DB_PATH: "/app/apps/api/data/kaneo.db"
      DISABLE_REGISTRATION: "` + disableRegistrationValue + `"
    volumes:
      - sqlite_data:/app/apps/api/data
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.backend.rule=Host(\"api.` + apiDomain + `\")"
      - "traefik.http.routers.backend.entrypoints=web"
`
	if useHttps {
		content += `      - "traefik.http.routers.backend.middlewares=backend-https-redirect"
      - "traefik.http.middlewares.backend-https-redirect.redirectscheme.scheme=https"
      - "traefik.http.routers.backend-secure.rule=Host(\"api.` + apiDomain + `\")"
      - "traefik.http.routers.backend-secure.entrypoints=websecure"
      - "traefik.http.routers.backend-secure.tls=true"
      - "traefik.http.routers.backend-secure.tls.certresolver=myresolver"
`
	}

	content += `      - "traefik.http.services.backend.loadbalancer.server.port=1337"
    networks:
      - traefik-net
    restart: unless-stopped

  frontend:
    image: ghcr.io/usekaneo/web:latest
    environment:
      KANEO_API_URL: "` + protocol + `://api.` + apiDomain + `"
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.frontend.rule=Host(\"` + apiDomain + `\")"
      - "traefik.http.routers.frontend.entrypoints=web"
`
	if useHttps {
		content += `      - "traefik.http.routers.frontend.middlewares=frontend-https-redirect"
      - "traefik.http.middlewares.frontend-https-redirect.redirectscheme.scheme=https"
      - "traefik.http.routers.frontend-secure.rule=Host(\"` + apiDomain + `\")"
      - "traefik.http.routers.frontend-secure.entrypoints=websecure"
      - "traefik.http.routers.frontend-secure.tls=true"
      - "traefik.http.routers.frontend-secure.tls.certresolver=myresolver"
`
	}

	content += `      - "traefik.http.services.frontend.loadbalancer.server.port=80"
    networks:
      - traefik-net
    restart: unless-stopped

networks:
  traefik-net:
    driver: bridge

volumes:
  sqlite_data:
`
	if useHttps {
		content += `  letsencrypt:
`
	}

	writeToFile("compose.yml", content)
	fmt.Printf("%s Created Traefik Docker Compose file: %s\n", success("✓"), bold("compose.yml"))
	fmt.Printf("%s After deployment, Kaneo will be available at: %s://%s\n", info("ℹ"), protocol, apiDomain)
	fmt.Printf("%s The API will be available at: %s://api.%s\n", info("ℹ"), protocol, apiDomain)
	fmt.Println("Make sure your DNS settings point both domains to your server's IP address.")
}

func createNginxComposeFile() {
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")

	protocol := "http"
	if useHttps {
		protocol = "https"
	}

	disableRegistrationValue := "false"
	if disableRegister {
		disableRegistrationValue = "true"
	}

	content := `services:
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/conf.d/default.conf
    depends_on:
      - backend
      - frontend
    restart: unless-stopped

  backend:
    image: ghcr.io/usekaneo/api:latest
    environment:
      JWT_ACCESS: "` + jwtToken + `"
      DB_PATH: "/app/apps/api/data/kaneo.db"
      DISABLE_REGISTRATION: "` + disableRegistrationValue + `"
    expose:
      - "1337"
    restart: unless-stopped
    volumes:
      - sqlite_data:/app/apps/api/data

  frontend:
    image: ghcr.io/usekaneo/web:latest
    environment:
      KANEO_API_URL: "` + protocol + `://` + domain + `/api"
    expose:
      - "5173"
    restart: unless-stopped

volumes:
  sqlite_data:
`
	writeToFile("compose.yml", content)
	fmt.Printf("%s Created Nginx Docker Compose file: %s\n", success("✓"), bold("compose.yml"))

	nginxContent := `server {
    listen 80;
    server_name ` + domain + `;

    # Security headers
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;

    # Frontend proxy
    location / {
        proxy_pass http://frontend:5173;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_cache_bypass $http_upgrade;
    }

    # API proxy with proper path rewriting
    location /api/ {
        proxy_pass http://backend:1337/;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_buffering off;
        proxy_request_buffering off;
    }

    # Deny access to hidden files
    location ~ /\\. {
        deny all;
        access_log off;
        log_not_found off;
    }
}
`
	writeToFile("nginx.conf", nginxContent)
	fmt.Printf("%s Created Nginx configuration file: %s\n", success("✓"), bold("nginx.conf"))

	if useHttps {
		fmt.Printf("%s HTTPS is enabled. You will need to configure SSL certificates for Nginx.\n", info("ℹ"))
		fmt.Printf("%s You can use Certbot to obtain Let's Encrypt certificates:\n", info("ℹ"))
		fmt.Printf("  %s\n", bold("sudo certbot --nginx -d "+domain))
	}

	fmt.Printf("%s After deployment, Kaneo will be available at: %s://%s\n", info("ℹ"), protocol, domain)
}

func writeToFile(filename, content string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("%s Error creating file %s: %v\n", errorC("✗"), filename, err)
		os.Exit(1)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		fmt.Printf("%s Error writing to file %s: %v\n", errorC("✗"), filename, err)
		os.Exit(1)
	}
}

func startContainers() {
	fmt.Printf("%s Do you want to start the containers now? (y/n) ", info("?"))
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response == "y" || response == "yes" {
		fmt.Printf("%s Starting containers...\n", info("ℹ"))
		cmd := exec.Command("docker", "compose", "up", "-d")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("%s Error starting containers: %v\n", errorC("✗"), err)
			return
		}
		fmt.Printf("%s Containers started successfully!\n", success("✓"))
	} else {
		fmt.Printf("%s You can start the containers later with: %s\n", info("ℹ"), bold("docker compose up -d"))
	}
}
