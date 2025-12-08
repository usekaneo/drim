# drim

One-click deployment tool for self-hosted Kaneo instances.

## Installation

Download the latest binary for your platform:

```bash
# Linux
curl -L https://github.com/usekaneo/drim/releases/latest/download/drim_Linux_x86_64 -o drim
chmod +x drim
sudo mv drim /usr/local/bin/

# macOS (Apple Silicon)
curl -L https://github.com/usekaneo/drim/releases/latest/download/drim_Darwin_arm64 -o drim
chmod +x drim
sudo mv drim /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/usekaneo/drim/releases/latest/download/drim_Darwin_x86_64 -o drim
chmod +x drim
sudo mv drim /usr/local/bin/
```

## Quick Start

Deploy Kaneo with a single command:

```bash
drim setup
```

You'll be prompted for your domain (optional). If you skip it, Kaneo will be accessible at `http://localhost`.

For production with HTTPS, provide your domain when prompted. Make sure your domain's DNS A record points to your server.

## Commands

- `drim setup` - Deploy Kaneo
- `drim start` - Start services
- `drim stop` - Stop services
- `drim restart` - Restart services
- `drim upgrade` - Update to latest version
- `drim configure` - Edit configuration
- `drim uninstall` - Remove Kaneo

## What Gets Installed

- PostgreSQL database
- Kaneo API (`ghcr.io/usekaneo/api:latest`)
- Kaneo Web (`ghcr.io/usekaneo/web:latest`)
- Caddy reverse proxy with automatic HTTPS

## Configuration

After setup, you can configure Kaneo by editing the `.env` file:

```bash
drim configure
```

See the [official Kaneo documentation](https://kaneo.app/docs/installation/environment-variables) for all available environment variables.

### Optional Features

Edit `.env` and uncomment variables to enable:

- **GitHub OAuth**: Set `GITHUB_CLIENT_ID` and `GITHUB_CLIENT_SECRET`
- **Google OAuth**: Set `GOOGLE_CLIENT_ID` and `GOOGLE_CLIENT_SECRET`
- **Email Auth**: Configure SMTP settings

## Requirements

- Docker 20.10 or later
- Docker Compose V2
- Linux, macOS, or Windows with WSL

drim will attempt to install Docker automatically on supported Linux distributions.

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

## Links

- [Kaneo Documentation](https://kaneo.app/docs)
- [Report Issues](https://github.com/usekaneo/drim/issues)
