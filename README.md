# Go Clean Architecture SaaS Starter Kit

A production-ready SaaS starter kit built with Go, following Clean Architecture principles. This template provides essential SaaS features including JWT authentication, multi-tenancy, subscription management, and more.

## 🚀 Features

- **Clean Architecture**: Separation of concerns with clear boundaries between layers
- **Flexible Configuration**: Support for .env, config.json, or environment variables with priority override
- **JWT Authentication**: Access tokens (1 hour) + Refresh tokens (7 days)
- **Multi-Tenancy**: Organization-first design with role-based access control
- **UUID Primary Keys**: CHAR(36) format for global uniqueness and security
- **Soft Delete**: Data retention with deleted_at timestamps on all tables
- **Subscription Management**: Tiered plans with upgrade/downgrade support
- **Health Checks**: Liveness and readiness probes for monitoring
- **Database Migrations**: Version-controlled schema changes with golang-migrate
- **Comprehensive Testing**: 59 passing tests covering all features and database schema
- **Docker Support**: Multi-stage builds with docker-compose
- **Frontend Friendly**: CORS enabled, rate limiting disabled by default
- **Audit Logging**: Optional audit trail for compliance (table created, implementation optional)

## 📋 Tech Stack

- **Language**: Go 1.23+
- **Web Framework**: Fiber v2
- **ORM**: GORM
- **Database**: MySQL 8.0
- **Authentication**: JWT with golang-jwt/jwt/v5
- **Configuration**: Viper (supports .env and config.json)
- **Logging**: Logrus
- **Validation**: go-playground/validator/v10
- **Password Hashing**: bcrypt

## 🏗️ Architecture

```
cmd/web/              - Application entry point
internal/
  ├── config/         - Configuration and bootstrap
  ├── entity/         - Domain entities (database models)
  ├── model/          - Request/Response DTOs
  │   └── converter/  - Entity to DTO converters
  ├── repository/     - Data access layer
  ├── usecase/        - Business logic layer
  └── delivery/       - Presentation layer
      └── http/       - HTTP handlers and routes
          └── middleware/ - Auth middleware
pkg/
  └── jwt/            - JWT service
db/migrations/        - Database migration files
```

## 🎯 API Endpoints

### Health Checks
- `GET /health` - Health check
- `GET /ready` - Readiness check (includes DB connection test)

### Authentication (Public)
- `POST /api/v1/auth/register` - Register new organization + user
- `POST /api/v1/auth/login` - Login with email/password
- `POST /api/v1/auth/refresh` - Refresh access token

### Authentication (Protected)
- `DELETE /api/v1/auth/logout` - Logout (clears refresh token)

### Users (Protected)
- `GET /api/v1/users/current` - Get current user
- `PATCH /api/v1/users/current` - Update current user

### Organizations (Protected)
- `GET /api/v1/organizations/current` - Get current organization
- `PATCH /api/v1/organizations/current` - Update organization
- `GET /api/v1/organizations/members` - List organization members
- `DELETE /api/v1/organizations/members/:userId` - Remove member

### Subscriptions (Protected)
- `GET /api/v1/subscriptions/current` - Get current subscription
- `POST /api/v1/subscriptions/upgrade` - Upgrade/downgrade plan
- `POST /api/v1/subscriptions/cancel` - Cancel subscription

## 🛠️ Quick Start

### Prerequisites

- Go 1.23 or higher
- MySQL 8.0 or higher
- Make (optional, but recommended)

### 1. Clone and Setup

```bash
# Clone the repository
git clone <your-repo-url>
cd go-clean-arch-saas

# Choose your configuration method:

# Option A: Using .env file (recommended for Docker)
cp .env.example .env
nano .env  # Edit database credentials and JWT secret

# Option B: Using config.json (already included)
nano config.json  # Edit configuration directly

# Option C: Mix both (env vars override config.json)
cp .env.example .env
# Edit only specific values you want to override
```

### 2. Database Setup

```bash
# Create database
mysql -u root -p -e "CREATE DATABASE go_clean_arch_saas"

# Run migrations
make migrate-up

# (Optional) Seed with demo data
make seed
```

### 3. Run Application

**Development mode with hot reload:**
```bash
make dev
```

**Production mode:**
```bash
make build
./bin/app
```

**Using Docker:**
```bash
make docker-up
```

## 📝 Configuration

The application supports **flexible configuration** - choose your preferred method:

### Configuration Priority (Highest to Lowest)
1. **Environment Variables** - Direct OS environment variables
2. **`.env` File** - For Docker/containerized environments
3. **`config.json` File** - Traditional JSON configuration
4. **Default Values** - Built-in fallbacks

### Option 1: Using .env File

Copy `.env.example` to `.env` and customize:

```bash
cp .env.example .env
```

```env
# Application
APP_NAME=go-clean-arch-saas
APP_ENV=development

# Server
WEB_PORT=3000
WEB_PREFORK=false

# Database
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=
DB_NAME=go_clean_arch_saas
DB_POOL_IDLE=10
DB_POOL_MAX=100
DB_POOL_LIFETIME=300

# JWT (minimum 32 characters recommended)
JWT_SECRET=your-secret-key-change-in-production-min-32-chars
JWT_ACCESS_EXPIRE_MINUTES=60
JWT_REFRESH_EXPIRE_DAYS=7

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
CORS_ALLOWED_METHODS=GET,POST,PUT,PATCH,DELETE
CORS_ALLOWED_HEADERS=Origin,Content-Type,Accept,Authorization

# Rate Limiting (disabled by default)
RATE_LIMIT_ENABLED=false
RATE_LIMIT_RPM=1000

# Logging (6=Trace, 5=Debug, 4=Info, 3=Warn, 2=Error, 1=Fatal, 0=Panic)
LOG_LEVEL=6
```

### Option 2: Using config.json File

Edit `config.json` directly:

```json
{
  "app": {
    "name": "go-clean-arch-saas",
    "env": "development"
  },
  "web": {
    "prefork": false,
    "port": 3000
  },
  "database": {
    "username": "root",
    "password": "",
    "host": "localhost",
    "port": 3306,
    "name": "go_clean_arch_saas",
    "pool": {
      "idle": 10,
      "max": 100,
      "lifetime": 300
    }
  },
  "jwt": {
    "secret": "your-secret-key-change-in-production-min-32-chars",
    "access_expire_minutes": 60,
    "refresh_expire_days": 7
  },
  "cors": {
    "allowed_origins": "http://localhost:3000,http://localhost:8080",
    "allowed_methods": "GET,POST,PUT,PATCH,DELETE",
    "allowed_headers": "Origin,Content-Type,Accept,Authorization"
  },
  "rate_limit": {
    "enabled": false,
    "rpm": 1000
  },
  "log": {
    "level": 6
  }
}
```

### Mix Both Methods

You can use `config.json` for base configuration and override specific values with `.env`:

```bash
# config.json has DB_NAME=go_clean_arch_saas
# .env overrides it:
DB_NAME=my_custom_database
```

### Configuration Keys Reference

| Environment Variable | JSON Path | Description | Default |
|---------------------|-----------|-------------|---------|
| `APP_NAME` | `app.name` | Application name | `go-clean-arch-saas` |
| `APP_ENV` | `app.env` | Environment | `development` |
| `WEB_PORT` | `web.port` | HTTP port | `3000` |
| `WEB_PREFORK` | `web.prefork` | Enable prefork mode | `false` |
| `DB_HOST` | `database.host` | Database host | `localhost` |
| `DB_PORT` | `database.port` | Database port | `3306` |
| `DB_USERNAME` | `database.username` | Database user | `root` |
| `DB_PASSWORD` | `database.password` | Database password | `` |
| `DB_NAME` | `database.name` | Database name | `go_clean_arch_saas` |
| `DB_POOL_IDLE` | `database.pool.idle` | Idle connections | `10` |
| `DB_POOL_MAX` | `database.pool.max` | Max connections | `100` |
| `DB_POOL_LIFETIME` | `database.pool.lifetime` | Connection lifetime (seconds) | `300` |
| `JWT_SECRET` | `jwt.secret` | JWT signing secret | - |
| `JWT_ACCESS_EXPIRE_MINUTES` | `jwt.access_expire_minutes` | Access token expiry | `60` |
| `JWT_REFRESH_EXPIRE_DAYS` | `jwt.refresh_expire_days` | Refresh token expiry | `7` |
| `CORS_ALLOWED_ORIGINS` | `cors.allowed_origins` | CORS origins | `http://localhost:3000,http://localhost:8080` |
| `CORS_ALLOWED_METHODS` | `cors.allowed_methods` | CORS methods | `GET,POST,PUT,PATCH,DELETE` |
| `CORS_ALLOWED_HEADERS` | `cors.allowed_headers` | CORS headers | `Origin,Content-Type,Accept,Authorization` |
| `RATE_LIMIT_ENABLED` | `rate_limit.enabled` | Enable rate limiting | `false` |
| `RATE_LIMIT_RPM` | `rate_limit.rpm` | Requests per minute | `1000` |
| `LOG_LEVEL` | `log.level` | Log level (0-6) | `6` |

## 🗄️ Database Schema

### Core Tables

- **organizations** - Tenant/organization data
- **users** - User accounts with organization relation
- **organization_members** - User roles within organizations
- **plans** - Subscription plan definitions
- **subscriptions** - Active organization subscriptions
- **audit_logs** - Optional audit trail (table ready, logging not implemented)

### UUID Primary Keys

All tables use **UUID (CHAR(36))** as primary keys for:
- ✅ Global uniqueness across distributed systems
- ✅ Security (non-sequential IDs prevent enumeration attacks)
- ✅ Client-side ID generation capability
- ✅ Easier data merging from multiple sources

Example ID format: `550e8400-e29b-41d4-a716-446655440000`

### Soft Delete Support

All tables implement **soft delete** with `deleted_at` column (BIGINT, Unix timestamp in milliseconds):
- ✅ Data retention for audit and compliance
- ✅ Accidental deletion recovery
- ✅ Historical data preservation
- ✅ Referential integrity maintained

```sql
-- Active records (not deleted)
WHERE deleted_at IS NULL

-- Soft deleted records
WHERE deleted_at IS NOT NULL

-- Restore soft deleted record
UPDATE users SET deleted_at = NULL WHERE id = ?
```

Each `deleted_at` column has an index (e.g., `idx_users_deleted`) for optimal query performance.

### Default Plans

The seed script creates three plans:
- **Free**: $0/month - Basic features for testing
- **Pro**: $29/month - Advanced features for growing teams
- **Enterprise**: $99/month - Full features with priority support

## 🔐 Authentication Flow

### Registration
1. User submits name, email, password, organization_name
2. System creates organization with unique slug
3. System creates user with hashed password
4. System adds user as organization owner
5. System creates free subscription
6. Returns JWT access token + refresh token

### Login
1. User submits email and password
2. System verifies credentials
3. System generates access token (1 hour expiry)
4. System generates refresh token (7 days expiry) and stores in DB
5. Returns both tokens

### Token Refresh
1. Client submits refresh token
2. System validates token expiry from database
3. System generates new access token
4. Returns new access token (refresh token unchanged)

### Protected Routes
- All requests must include `Authorization: Bearer <access_token>` header
- Middleware validates JWT signature and expiry
- Middleware injects user_id, email, organization_id into context

## 🧪 Testing

The project includes comprehensive unit tests covering all API endpoints and database features.

**Test Coverage:**
- ✅ Health checks (2 tests)
- ✅ Authentication (10 tests) - register, login, refresh, logout
- ✅ User management (8 tests) - get, update
- ✅ Organization management (10 tests) - get, update, members
- ✅ Subscription management (11 tests) - get, upgrade, cancel, workflow
- ✅ Database schema (18 tests) - UUID validation, soft delete functionality

**Total: 59 passing tests**

```bash
# Run all tests (make sure database is created and migrations are applied)
make test

# Run tests with coverage
make test-coverage

# Run specific test
go test -v ./test/ -run TestRegister
go test -v ./test/ -run TestGetCurrentUser
go test -v ./test/ -run TestSoftDelete
```

### Soft Delete Testing

The test suite validates soft delete functionality:
- Manual soft delete operations
- Query filtering with `deleted_at IS NULL`
- Timestamp verification
- Recovery capability

Example from `test/schema_test.go`:
```go
// Manual soft delete
now := time.Now().UnixMilli()
db.Model(&entity.User{}).Where("id = ?", userID).Update("deleted_at", now)

// Verify filtering
db.Where("id = ? AND deleted_at IS NULL", userID).First(&user)
// Should not find soft-deleted record
```

## 📦 Available Make Commands

```bash
make setup          # Initial project setup
make dev            # Run with hot reload (requires air)
make build          # Build binary
make run            # Run built binary
make test           # Run tests
make test-coverage  # Run tests with coverage report
make migrate-up     # Run migrations
make migrate-down   # Rollback migrations
make migrate-create # Create new migration
make seed           # Seed database with demo data
make clean          # Clean build artifacts
make docker-build   # Build Docker image
make docker-up      # Start with docker-compose
make docker-down    # Stop docker containers
```

## 🐳 Docker Deployment

The project includes a multi-stage Dockerfile and docker-compose configuration:

```bash
# Build and start services
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop services
docker-compose down
```

## 🔧 Customization Guide

### Adding New Entities

1. Create entity in `internal/entity/`
2. Create migration in `db/migrations/`
3. Create model DTOs in `internal/model/`
4. Create converter in `internal/model/converter/`
5. Create repository in `internal/repository/`
6. Create usecase in `internal/usecase/`
7. Create controller in `internal/delivery/http/`
8. Register in `internal/config/app.go` Bootstrap
9. Add routes in `internal/delivery/http/route/route.go`

### Implementing Audit Logging

The `audit_logs` table is created but not actively used. To implement:

1. Create an audit service in `pkg/audit/`
2. Call audit service from usecases after important actions
3. Example: Log user creation, organization updates, subscription changes

## 🔥 Quick Test

```bash
# Start the server
make dev

# Register a new user
curl -X POST http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123",
    "organization_name": "Acme Corp"
  }'

# Login
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'

# Use the access_token from login response
curl -X GET http://localhost:3000/api/v1/users/current \
  -H "Authorization: Bearer <your-access-token>"
```

## 📚 Additional Documentation

- **[Configuration Guide](docs/CONFIGURATION.md)** - Comprehensive configuration documentation with examples
- **[API Specification](api/api-spec.json)** - OpenAPI/Swagger documentation *(to be updated)*
- **[Testing Guide](docs/TESTING.md)** - Full testing documentation *(coming soon)*

## 🤝 Contributing

This is a starter template. Fork it and make it your own!

## 📄 License

See [LICENSE.txt](LICENSE.txt) for details.

## 💡 Notes

- **No Example Business Logic**: This template is intentionally clean - no contact management or other example features
- **Frontend Friendly**: Rate limiting is disabled by default to work seamlessly with frontend frameworks (Next.js, React, Vue, Nuxt)
- **Production Ready**: Includes health checks, migrations, Docker support, and proper error handling
- **Multi-Tenancy**: Organization-first design ensures proper data isolation
- **Scalable**: Clean architecture makes it easy to add features without technical debt

---

**Built with ❤️ using Go Clean Architecture**