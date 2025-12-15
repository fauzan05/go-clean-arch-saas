# Go Clean Architecture SaaS Starter Kit

A production-ready SaaS starter kit built with Go, following Clean Architecture principles. This template provides essential SaaS features including JWT authentication, multi-tenancy, subscription management, and more.

## 🚀 Features

- **Clean Architecture**: Separation of concerns with clear boundaries between layers
- **JWT Authentication**: Access tokens (1 hour) + Refresh tokens (7 days)
- **Multi-Tenancy**: Organization-first design with role-based access control
- **Subscription Management**: Tiered plans with upgrade/downgrade support
- **Health Checks**: Liveness and readiness probes for monitoring
- **Database Migrations**: Version-controlled schema changes
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

# Copy environment file
cp .env.example .env

# Edit .env with your database credentials and JWT secret
nano .env
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

Configuration can be provided via `.env` file or `config.json`.

### Environment Variables

```env
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=go_clean_arch_saas

# JWT
JWT_SECRET=your-super-secret-key-change-in-production
JWT_ACCESS_EXPIRE_MINUTES=60
JWT_REFRESH_EXPIRE_DAYS=7

# Server
PORT=3000
PREFORK=false

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
CORS_ALLOWED_METHODS=GET,POST,PUT,PATCH,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Accept,Authorization,Content-Type,X-CSRF-Token

# Rate Limiting (disabled by default for frontend-friendly setup)
RATE_LIMIT_ENABLED=false
RATE_LIMIT_MAX=20
RATE_LIMIT_WINDOW=60s

# Logging
LOG_LEVEL=6  # 6=Trace, 5=Debug, 4=Info, 3=Warn, 2=Error, 1=Fatal, 0=Panic
```

## 🗄️ Database Schema

### Core Tables

- **organizations** - Tenant/organization data
- **users** - User accounts with organization relation
- **organization_members** - User roles within organizations
- **plans** - Subscription plan definitions
- **subscriptions** - Active organization subscriptions
- **audit_logs** - Optional audit trail (table ready, logging not implemented)

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

The project includes comprehensive unit tests covering all API endpoints.

**Test Coverage:**
- ✅ Health checks (2 tests)
- ✅ Authentication (10 tests) - register, login, refresh, logout
- ✅ User management (8 tests) - get, update
- ✅ Organization management (10 tests) - get, update, members
- ✅ Subscription management (11 tests) - get, upgrade, cancel, workflow

```bash
# Run all tests (make sure database is created and migrations are applied)
make test

# Run tests with coverage
make test-coverage

# Run specific test
go test -v ./test/ -run TestRegister
go test -v ./test/ -run TestGetCurrentUser
```

📖 **Full testing guide**: See [docs/TESTING.md](docs/TESTING.md) for detailed information

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

- [API Specification](api/api-spec.json) - OpenAPI/Swagger documentation *(to be updated)*
- More detailed docs coming soon

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