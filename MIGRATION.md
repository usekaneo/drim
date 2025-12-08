# Migrating Existing Kaneo to drim

This guide helps you migrate an existing Kaneo installation to drim without losing data.

## Before You Start

**Important:** drim uses Caddy instead of nginx. Your data will be preserved, but the reverse proxy configuration will change.

### What You Need

- SSH access to your server
- Sudo privileges
- Your existing Kaneo installation directory
- Backup of your current setup (recommended)

## Migration Steps

### Step 1: Backup Everything

```bash
# Navigate to your existing Kaneo directory
cd /path/to/your/kaneo

# Create backup directory
mkdir -p ~/kaneo-backup-$(date +%Y%m%d)
cd ~/kaneo-backup-$(date +%Y%m%d)

# Backup configuration files
cp /path/to/your/kaneo/docker-compose.yml ./
cp /path/to/your/kaneo/.env ./

# Backup nginx configuration (if separate)
cp /etc/nginx/sites-available/kaneo ./nginx-config 2>/dev/null || true

# Backup database
docker compose exec -T postgres pg_dump -U kaneo kaneo > database-backup.sql

# List volumes to document them
docker volume ls | grep kaneo > volumes.txt
```

### Step 2: Export Current Environment Variables

```bash
# Save your current environment variables
cd /path/to/your/kaneo
cat .env > ~/kaneo-backup-$(date +%Y%m%d)/env-backup.txt
```

### Step 3: Stop Current Services

```bash
cd /path/to/your/kaneo
docker compose stop
```

**Note:** We're stopping, not removing. Your volumes are safe.

### Step 4: Identify Your Volume Names

```bash
# List your volumes
docker volume ls | grep postgres
docker volume ls | grep kaneo

# Inspect to confirm data location
docker volume inspect <volume-name>
```

Example output:
```
kaneo_postgres_data
kaneo_api_data
```

### Step 5: Install drim

```bash
curl -fsSL https://assets.kaneo.app/install.sh | sh
```

### Step 6: Move to New Directory

```bash
# Create a new directory for drim-managed Kaneo
mkdir -p ~/kaneo-drim
cd ~/kaneo-drim
```

### Step 7: Run drim Setup

```bash
drim setup
```

When prompted for domain, enter your existing domain.

### Step 8: Stop drim Services

```bash
drim stop
```

### Step 9: Update docker-compose.yml to Use Existing Volumes

Edit the generated `docker-compose.yml`:

```bash
nano docker-compose.yml
```

Update the postgres volume to use your existing volume:

```yaml
services:
  postgres:
    image: postgres:16-alpine
    env_file: .env
    volumes:
      - postgres_data:/var/lib/postgresql/data  # This line
    # ... rest of config

volumes:
  postgres_data:
    external: true
    name: kaneo_postgres_data  # Your existing volume name
  caddy_data:
  caddy_config:
```

### Step 10: Update .env with Your Existing Credentials

```bash
nano .env
```

Replace the generated values with your existing ones:

```env
KANEO_CLIENT_URL=https://your-domain.com
KANEO_API_URL=https://your-domain.com/api

DATABASE_URL=postgresql://your_user:your_password@postgres:5432/your_db
POSTGRES_DB=your_db
POSTGRES_USER=your_user
POSTGRES_PASSWORD=your_existing_password

AUTH_SECRET=your_existing_secret

DOMAIN=your-domain.com
```

### Step 11: Update Nginx (If You Want to Keep It)

If you prefer to keep nginx instead of Caddy:

**Option A: Remove Caddy and Use Existing Nginx**

Edit `docker-compose.yml` and remove the caddy service:

```yaml
# Remove or comment out:
# caddy:
#   image: caddy:2-alpine
#   ...
```

Then configure your existing nginx to proxy to the new containers:

```nginx
# /etc/nginx/sites-available/kaneo
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl;
    server_name your-domain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location /api {
        proxy_pass http://localhost:1337;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location / {
        proxy_pass http://localhost:5173;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

Update the docker-compose.yml to expose ports:

```yaml
services:
  api:
    image: ghcr.io/usekaneo/api:latest
    env_file: .env
    ports:
      - "1337:1337"
    # ... rest

  web:
    image: ghcr.io/usekaneo/web:latest
    env_file: .env
    ports:
      - "5173:5173"
    # ... rest
```

**Option B: Switch to Caddy (Recommended)**

If you want to use Caddy (automatic HTTPS), update your DNS and stop nginx:

```bash
# Stop nginx
sudo systemctl stop nginx
sudo systemctl disable nginx

# Start with Caddy
drim start
```

### Step 12: Start Services

```bash
drim start
```

### Step 13: Verify Everything Works

```bash
# Check services are running
docker compose ps

# Check logs
docker compose logs -f

# Test database connection
docker compose exec postgres psql -U kaneo -d kaneo -c "SELECT COUNT(*) FROM users;"

# Test web access
curl https://your-domain.com
```

## Rollback Plan

If something goes wrong:

```bash
# Stop drim services
cd ~/kaneo-drim
drim stop

# Go back to old directory
cd /path/to/your/kaneo

# Start old services
docker compose up -d

# Restore nginx if you stopped it
sudo systemctl start nginx
```

## Common Issues

### Issue: Database Connection Failed

**Solution:** Check that the database credentials in `.env` match your existing setup.

```bash
# View your backup credentials
cat ~/kaneo-backup-*/env-backup.txt | grep POSTGRES
```

### Issue: Volumes Not Found

**Solution:** Ensure you're using the correct volume name.

```bash
# List all volumes
docker volume ls

# Update docker-compose.yml with exact name
```

### Issue: Port Already in Use

**Solution:** Another service is using port 80 or 443.

```bash
# Check what's using the port
sudo lsof -i :80
sudo lsof -i :443

# Stop the conflicting service or use nginx option
```

### Issue: Lost Environment Variables

**Solution:** Restore from backup.

```bash
# Copy your backed up .env
cp ~/kaneo-backup-*/env-backup.txt ~/kaneo-drim/.env
drim restart
```

## After Migration

Once everything is working:

### Clean Up Old Setup (Optional)

```bash
# Remove old containers (not volumes!)
cd /path/to/your/kaneo
docker compose down

# Keep volumes for safety
# Only remove after confirming everything works for a week
```

### Update DNS (If Changed)

If you moved servers, update your DNS A record to point to the new server IP.

### Monitor Logs

```bash
cd ~/kaneo-drim
docker compose logs -f
```

### Set Up Backups

```bash
# Add to crontab
crontab -e

# Daily backup at 2 AM
0 2 * * * cd ~/kaneo-drim && docker compose exec -T postgres pg_dump -U kaneo kaneo > ~/backups/kaneo-$(date +\%Y\%m\%d).sql
```

## Need Help?

- Check logs: `docker compose logs -f`
- Verify volumes: `docker volume ls`
- Test database: `docker compose exec postgres psql -U kaneo -d kaneo`
- [Open an issue](https://github.com/usekaneo/drim/issues)

