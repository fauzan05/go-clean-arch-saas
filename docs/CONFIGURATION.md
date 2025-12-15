# Configuration Guide

This application provides flexible configuration management with support for multiple sources.

## Configuration Priority

The application loads configuration in the following priority order (highest to lowest):

1. **Environment Variables** - Direct OS environment variables (e.g., `export DB_HOST=localhost`)
2. **`.env` File** - Environment file in the root directory
3. **`config.json` File** - JSON configuration file in the root directory
4. **Default Values** - Built-in fallback values

This means:
- Environment variables override everything
- `.env` values override `config.json`
- `config.json` provides base configuration
- Defaults are used if nothing else is set

## Quick Start Examples

### Scenario 1: Simple Development Setup

Use the provided `config.json` as-is:

```bash
# No setup needed! Just run:
make dev
```

### Scenario 2: Docker/Container Setup

Use `.env` file for easy environment variable management:

```bash
cp .env.example .env
# Edit .env with your settings
docker-compose up
```

### Scenario 3: Override Single Values

Keep `config.json` but override database name:

```bash
export DB_NAME=my_custom_database
make run
```

### Scenario 4: Production with Secrets

Use `config.json` for non-sensitive data, environment variables for secrets:

```bash
# config.json has all settings
# Set sensitive values via environment:
export DB_PASSWORD=super_secure_password
export JWT_SECRET=production-secret-key-minimum-32-characters-long
./bin/app
```

## Configuration Reference

### Application Settings

| Key | Env Var | Description | Default |
|-----|---------|-------------|---------|
| `app.name` | `APP_NAME` | Application name shown in logs | `go-clean-arch-saas` |
| `app.env` | `APP_ENV` | Environment (development, staging, production) | `development` |

### Web Server Settings

| Key | Env Var | Description | Default |
|-----|---------|-------------|---------|
| `web.port` | `WEB_PORT` | HTTP server port | `3000` |
| `web.prefork` | `WEB_PREFORK` | Enable Fiber prefork mode for better performance | `false` |

**Prefork Mode**: Set to `true` in production for better performance with multiple CPU cores.

### Database Settings

| Key | Env Var | Description | Default |
|-----|---------|-------------|---------|
| `database.host` | `DB_HOST` | MySQL host address | `localhost` |
| `database.port` | `DB_PORT` | MySQL port | `3306` |
| `database.username` | `DB_USERNAME` | MySQL username | `root` |
| `database.password` | `DB_PASSWORD` | MySQL password | `` (empty) |
| `database.name` | `DB_NAME` | Database name | `go_clean_arch_saas` |
| `database.pool.idle` | `DB_POOL_IDLE` | Minimum idle connections in pool | `10` |
| `database.pool.max` | `DB_POOL_MAX` | Maximum open connections | `100` |
| `database.pool.lifetime` | `DB_POOL_LIFETIME` | Connection lifetime in seconds | `300` |

**Connection Pool Tuning**:
- `idle`: Number of connections kept ready (higher = faster response, more memory)
- `max`: Maximum connections (should be less than MySQL's `max_connections`)
- `lifetime`: Recycle connections periodically to prevent stale connections

### JWT Settings

| Key | Env Var | Description | Default |
|-----|---------|-------------|---------|
| `jwt.secret` | `JWT_SECRET` | Secret key for signing JWT tokens | `your-secret-key-change-in-production-min-32-chars` |
| `jwt.access_expire_minutes` | `JWT_ACCESS_EXPIRE_MINUTES` | Access token expiration in minutes | `60` |
| `jwt.refresh_expire_days` | `JWT_REFRESH_EXPIRE_DAYS` | Refresh token expiration in days | `7` |

**Security Notes**:
- JWT secret MUST be changed in production
- Minimum 32 characters recommended for security
- Access tokens are short-lived (1 hour) for security
- Refresh tokens allow re-authentication without login (7 days)

### CORS Settings

| Key | Env Var | Description | Default |
|-----|---------|-------------|---------|
| `cors.allowed_origins` | `CORS_ALLOWED_ORIGINS` | Comma-separated list of allowed origins | `http://localhost:3000,http://localhost:8080` |
| `cors.allowed_methods` | `CORS_ALLOWED_METHODS` | Comma-separated HTTP methods | `GET,POST,PUT,PATCH,DELETE` |
| `cors.allowed_headers` | `CORS_ALLOWED_HEADERS` | Comma-separated allowed headers | `Origin,Content-Type,Accept,Authorization` |

**Frontend Integration**:
```env
# For Next.js on port 3000
CORS_ALLOWED_ORIGINS=http://localhost:3000

# For multiple frontends
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173,https://app.example.com
```

### Rate Limiting Settings

| Key | Env Var | Description | Default |
|-----|---------|-------------|---------|
| `rate_limit.enabled` | `RATE_LIMIT_ENABLED` | Enable rate limiting | `false` |
| `rate_limit.rpm` | `RATE_LIMIT_RPM` | Requests per minute per IP | `1000` |

**Note**: Rate limiting is disabled by default for easier frontend development.

### Logging Settings

| Key | Env Var | Description | Default |
|-----|---------|-------------|---------|
| `log.level` | `LOG_LEVEL` | Log verbosity level (0-6) | `6` |

**Log Levels**:
- `6` - Trace (most verbose, shows SQL queries)
- `5` - Debug (development debugging)
- `4` - Info (production recommended)
- `3` - Warn (warnings only)
- `2` - Error (errors only)
- `1` - Fatal (fatal errors)
- `0` - Panic (critical panics)

## Configuration Examples

### .env File Example

```env
# Full example with all options
APP_NAME=my-saas-app
APP_ENV=production

WEB_PORT=8080
WEB_PREFORK=true

DB_HOST=db.example.com
DB_PORT=3306
DB_USERNAME=saas_user
DB_PASSWORD=secure_password_here
DB_NAME=production_db
DB_POOL_IDLE=20
DB_POOL_MAX=200
DB_POOL_LIFETIME=600

JWT_SECRET=production-super-secret-key-min-32-chars-please-change-this
JWT_ACCESS_EXPIRE_MINUTES=30
JWT_REFRESH_EXPIRE_DAYS=14

CORS_ALLOWED_ORIGINS=https://app.example.com,https://admin.example.com
CORS_ALLOWED_METHODS=GET,POST,PUT,PATCH,DELETE
CORS_ALLOWED_HEADERS=Origin,Content-Type,Accept,Authorization,X-Request-ID

RATE_LIMIT_ENABLED=true
RATE_LIMIT_RPM=100

LOG_LEVEL=4
```

### config.json Example

```json
{
  "app": {
    "name": "my-saas-app",
    "env": "production"
  },
  "web": {
    "prefork": true,
    "port": 8080
  },
  "database": {
    "username": "saas_user",
    "password": "secure_password_here",
    "host": "db.example.com",
    "port": 3306,
    "name": "production_db",
    "pool": {
      "idle": 20,
      "max": 200,
      "lifetime": 600
    }
  },
  "jwt": {
    "secret": "production-super-secret-key-min-32-chars",
    "access_expire_minutes": 30,
    "refresh_expire_days": 14
  },
  "cors": {
    "allowed_origins": "https://app.example.com,https://admin.example.com",
    "allowed_methods": "GET,POST,PUT,PATCH,DELETE",
    "allowed_headers": "Origin,Content-Type,Accept,Authorization,X-Request-ID"
  },
  "rate_limit": {
    "enabled": true,
    "rpm": 100
  },
  "log": {
    "level": 4
  }
}
```

## Environment-Specific Configurations

### Development

```env
APP_ENV=development
WEB_PORT=3000
DB_NAME=go_clean_arch_saas_dev
LOG_LEVEL=6
RATE_LIMIT_ENABLED=false
```

### Staging

```env
APP_ENV=staging
WEB_PORT=3000
DB_NAME=go_clean_arch_saas_staging
LOG_LEVEL=5
RATE_LIMIT_ENABLED=true
RATE_LIMIT_RPM=500
```

### Production

```env
APP_ENV=production
WEB_PORT=80
WEB_PREFORK=true
DB_NAME=go_clean_arch_saas_prod
DB_POOL_MAX=200
LOG_LEVEL=4
RATE_LIMIT_ENABLED=true
RATE_LIMIT_RPM=100
JWT_ACCESS_EXPIRE_MINUTES=30
```

## Docker Compose Configuration

When using Docker Compose, pass environment variables via `docker-compose.yml`:

```yaml
version: '3.8'
services:
  app:
    build: .
    environment:
      - DB_HOST=mysql
      - DB_NAME=go_clean_arch_saas
      - JWT_SECRET=${JWT_SECRET}
    env_file:
      - .env
```

Or use a separate `.env` file that docker-compose will automatically load.

## Kubernetes Configuration

Use ConfigMaps for non-sensitive data and Secrets for passwords:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  DB_HOST: "mysql-service"
  DB_NAME: "go_clean_arch_saas"
  LOG_LEVEL: "4"
---
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
stringData:
  DB_PASSWORD: "secure_password"
  JWT_SECRET: "production-secret-key"
```

## Troubleshooting

### Configuration Not Loading

1. Check file exists: `ls -la .env config.json`
2. Check file permissions: `chmod 644 .env config.json`
3. Verify JSON syntax: `cat config.json | jq .`
4. Check environment variables: `printenv | grep DB_`

### Environment Variables Not Working

Make sure variables are exported:
```bash
# Wrong
DB_HOST=localhost ./bin/app

# Correct
export DB_HOST=localhost
./bin/app

# Or inline
DB_HOST=localhost DB_NAME=test ./bin/app
```

### Config Priority Not Working

Remember the priority order:
1. OS Environment Variables (highest)
2. `.env` file
3. `config.json` file
4. Defaults (lowest)

Check what's actually loaded:
```bash
# Add debug logging to see loaded config
LOG_LEVEL=6 ./bin/app
```

## Best Practices

1. **Never commit `.env`** - Add to `.gitignore`
2. **Use `.env.example`** - Document all required variables
3. **Separate secrets** - Keep sensitive data out of `config.json`
4. **Environment-specific configs** - Use different values per environment
5. **Validate on startup** - App will panic if required config is missing
6. **Document changes** - Update `.env.example` when adding new config
7. **Use meaningful defaults** - Sensible defaults for development
8. **Encrypt production secrets** - Use secret management tools in production

## Adding New Configuration

To add a new configuration option:

1. **Add to viper.go defaults**:
```go
config.SetDefault("myfeature.enabled", true)
```

2. **Add environment binding**:
```go
config.BindEnv("myfeature.enabled", "MYFEATURE_ENABLED")
```

3. **Update .env.example**:
```env
MYFEATURE_ENABLED=true
```

4. **Update config.json**:
```json
{
  "myfeature": {
    "enabled": true
  }
}
```

5. **Document in this file** - Add to the reference table above
