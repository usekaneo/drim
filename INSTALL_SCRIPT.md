# Installation Script

The `install.sh` script provides a one-line installation for drim.

## Usage

### Basic Installation

```bash
curl -fsSL https://drim.kaneo.app/install.sh | sh
```

### With Options

```bash
# Silent mode (no output)
curl -fsSL https://drim.kaneo.app/install.sh | sh -s -- --silent

# Auto-run setup after install
curl -fsSL https://drim.kaneo.app/install.sh | sh -s -- --setup

# Install and setup with domain
curl -fsSL https://drim.kaneo.app/install.sh | sh -s -- --setup --domain=kaneo.example.com

# Combine options
curl -fsSL https://drim.kaneo.app/install.sh | sh -s -- --silent --setup --domain=kaneo.example.com
```

## Options

- `--silent` - Suppress output messages
- `--setup` - Automatically run `drim setup` after installation
- `--domain=DOMAIN` - Pass domain to setup (requires `--setup`)

## Hosting on Cloudflare R2

To host this script on Cloudflare R2:

1. Create an R2 bucket (e.g., `drim-install`)

2. Upload `install.sh` to the bucket

3. Configure a custom domain (e.g., `drim.kaneo.app`)

4. Set the R2 bucket to public or use a presigned URL

5. Configure your DNS:
   ```
   CNAME drim.kaneo.app -> [your-r2-bucket-url]
   ```

6. Users can then install with:
   ```bash
   curl -fsSL https://drim.kaneo.app/install.sh | sh
   ```

## What It Does

1. Displays drim banner
2. Detects operating system and CPU architecture
3. Downloads the appropriate binary from GitHub releases
4. Installs to `/usr/local/bin` or `/bin` (with sudo if needed)
5. Optionally runs setup if `--setup` flag is provided
6. Shows success message with version

## Supported Platforms

- **Linux**: x86_64, arm64
- **macOS**: x86_64 (Intel), arm64 (Apple Silicon)

## Security

The script:
- Downloads from official GitHub releases
- Verifies HTTP response codes
- Uses temporary files with cleanup
- Only requires sudo for installation to system directories

