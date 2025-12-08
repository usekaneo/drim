# drim

One-click deployment tool for self-hosted Kaneo instances.

All you need. Nothing you don't.

## Installation

Install drim with a single command:

```bash
curl -fsSL https://assets.kaneo.app/install.sh | sh
```

Or download the binary manually from [releases](https://github.com/usekaneo/drim/releases/latest).

## Quick Start

Deploy Kaneo in seconds:

```bash
drim setup
```

That's it. Your Kaneo instance is now running.

### Local Development

```bash
drim setup
# Press Enter when prompted for domain
# Access at http://localhost
```

### Production Deployment

```bash
drim setup
# Enter your domain when prompted (e.g., kaneo.example.com)
# Access at https://your-domain.com (HTTPS automatic)
```

Make sure your domain's DNS A record points to your server before setup.

## Commands

```bash
drim setup        # Deploy Kaneo
drim start        # Start services
drim stop         # Stop services
drim restart      # Restart services
drim upgrade      # Update Kaneo to latest version
drim update       # Update drim CLI to latest version
drim configure    # Edit configuration
drim uninstall    # Remove Kaneo
```

## What Gets Installed

When you run `drim setup`, the following services are deployed:

- **PostgreSQL 16** - Database
- **Kaneo API** - Backend service
- **Kaneo Web** - Frontend interface
- **Caddy** - Reverse proxy with automatic HTTPS

All services run in Docker containers with proper networking and health checks.

## Configuration

### Edit Environment Variables

```bash
drim configure
```

This opens `.env` in your default editor. After saving, services are restarted automatically.

### Optional Features

Uncomment variables in `.env` to enable:

**GitHub Authentication**
```env
GITHUB_CLIENT_ID=your_client_id
GITHUB_CLIENT_SECRET=your_client_secret
```

**Google Authentication**
```env
GOOGLE_CLIENT_ID=your_client_id
GOOGLE_CLIENT_SECRET=your_client_secret
```

**Email Authentication (SMTP)**
```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-password
```

See [Kaneo documentation](https://kaneo.app/docs/installation/environment-variables) for all available options.

## Requirements

- Docker 20.10+
- Docker Compose V2
- 2GB RAM minimum
- 10GB disk space

**Supported Platforms:** Linux, macOS, Windows (WSL)

drim will attempt to install Docker automatically on supported Linux distributions.

## Examples

### Silent Installation

```bash
curl -fsSL https://assets.kaneo.app/install.sh | sh -s -- --silent
```

### Install and Setup in One Command

```bash
curl -fsSL https://assets.kaneo.app/install.sh | sh -s -- --setup --domain=kaneo.example.com
```

### Update Everything

```bash
drim update    # Update drim CLI
drim upgrade   # Update Kaneo
```

### Check Logs

```bash
docker compose logs -f
docker compose logs -f api
```

## Building from Source

```bash
git clone https://github.com/usekaneo/drim.git
cd drim
go build -o drim .
```

Build for all platforms:

```bash
make build-all
```

## License

MIT License. See [LICENSE](LICENSE) for details.

## Migrating from Existing Setup

If you have an existing Kaneo installation and want to migrate to drim without losing data, see the [Migration Guide](MIGRATION.md).

## Links

- [Kaneo](https://kaneo.app)
- [Kaneo Documentation](https://kaneo.app/docs)
- [Migration Guide](MIGRATION.md)
- [Report Issues](https://github.com/usekaneo/drim/issues)
