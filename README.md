# Drim - Kaneo Deployment CLI

Drim is a command-line interface tool that helps you deploy and manage Kaneo, an open-source project management platform. It simplifies the deployment process using Docker and Docker Compose.

## Features

- 🚀 Easy deployment with multiple proxy options (Traefik, Nginx, or none)
- 🔒 Automatic HTTPS configuration with Let's Encrypt
- 💾 Database backup and restore functionality
- 📊 Container status monitoring
- 📝 Log viewing capabilities
- 🔄 Easy updates to the latest version

## Prerequisites

- Docker
- Docker Compose
- Go 1.21 or later (for building from source)

## Installation

### Using Pre-built Binaries

Download the latest release from the [releases page](https://github.com/usekaneo/kaneo/releases).

### Building from Source

```bash
git clone https://github.com/usekaneo/kaneo.git
cd kaneo/drim
go build
```


## Usage

### Deploy Kaneo

```bash
# Basic deployment
drim deploy
```

```bash
# Deploy with specific options
drim deploy --domain example.com --proxy traefik --https
```

### Manage Deployment
```bash
# View deployment status
drim status
```

```bash
# View logs
drim logs
drim logs backend --follow
```

### Stop deployment

```bash
drim stop
```

### Backup and Restore

```bash
# Backup database
drim backup

# Restore database
drim restore --file ./backups/kaneo-backup-20240101-120000.db
```


## Configuration Options

- `--domain`: Your domain name (e.g., kaneo.example.com)
- `--proxy`: Proxy type (traefik, nginx, or none)
- `--https`: Enable HTTPS (with Traefik or Nginx)
- `--jwt`: Custom JWT token (auto-generated if not provided)
- `--disable-register`: Disable user registration

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[MIT](LICENSE)

