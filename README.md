# drim

One-click deployment tool for self-hosted Kaneo instances.

All you need. Nothing you don't.

![Demo](demo/demo.gif)

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
# Access at http://localhost:5173
```

The unified Kaneo container serves the web app and API through port `5173`.

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
- **Kaneo** (`ghcr.io/usekaneo/kaneo:latest`) - Unified web and API service on port `5173`
- **Caddy** - Reverse proxy with automatic HTTPS

The Kaneo image contains both the web app and API. All services run in Docker containers with a shared network and health checks.

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
GITHUB_OAUTH_CLIENT_ID=your_client_id
GITHUB_OAUTH_CLIENT_SECRET=your_client_secret
```

**Email Authentication (SMTP)**
```env
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=your-email@example.com
SMTP_PASSWORD=your-password
SMTP_FROM=your-email@example.com
```

**Redis Pub/Sub (optional)**
```env
REDIS_URL=redis://your-redis-host:6379
```

See Kaneo's [environment variable example](https://github.com/usekaneo/kaneo/blob/main/.env.sample) and [documentation](https://kaneo.app/docs/core) for required and optional settings.

## Requirements

- Docker 20.10+
- Docker Compose V2
- 2GB RAM minimum
- 10GB disk space
- A public DNS record for production deployments

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
docker compose logs -f kaneo
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
- [Kaneo Documentation](https://kaneo.app/docs/core)
- [Migration Guide](MIGRATION.md)
- [Report Issues](https://github.com/usekaneo/drim/issues)
