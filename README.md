# Go Microservice

A production-ready Go microservice template with comprehensive features for building scalable and maintainable APIs.

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [Database Migrations](#database-migrations)
- [Development Commands](#development-commands)
- [API Endpoints](#api-endpoints)
- [Architecture](#architecture)
- [Security](#security)
- [Testing](#testing)
- [Contributing](#contributing)
- [Acknowledgments](#acknowledgments)

## Features

### Core
- **Fiber Framework** - High-performance HTTP framework built on fasthttp
- **PostgreSQL** - Primary database with GORM ORM
- **Redis** - Rate limiting, login protection, token revocation
- **JWT Authentication** - Access/refresh token pattern with brute-force protection

### Architecture & Patterns
- **Clean Architecture** - Handlers → Services → Repositories
- **Repository Pattern** - Data access abstraction with GORM base repository
- **Manual Dependency Injection** - Explicit constructor wiring in `main.go`
- **Interface-Driven Design** - Interfaces defined where consumed (Go idiom)
- **Input Validation** - Struct tag validation with `go-playground/validator`

### Developer Experience
- **Docker Support** - Multi-stage build (~15MB production image)
- **Database Migrations** - Embedded SQL migrations with golang-migrate
- **Makefile** - Common development commands
- **Testing Suite** - Unit tests with testify, integration tests with build tags
- **Code Quality** - golangci-lint with 16 linters configured

### Production Ready
- **Error Handling** - Structured JSON error envelope with error type mapping
- **Rate Limiting** - Redis-backed with in-memory fallback, separate GET/write limits
- **Middleware Stack** - Request ID, timing, security headers, CORS, trusted hosts
- **Structured Logging** - zerolog with JSON (production) or pretty console (development)
- **Graceful Shutdown** - Signal handling with ordered resource cleanup
- **Nginx** - Reverse proxy configuration

## Project Structure

```
go-microservice/
├── cmd/
│   └── server/
│       └── main.go                  # Entry point, DI wiring, graceful shutdown
│
├── internal/
│   ├── api/
│   │   ├── router.go               # Route setup, middleware stack
│   │   ├── handlers/
│   │   │   ├── auth.go             # Register, login, refresh, me
│   │   │   ├── users.go            # User CRUD endpoints
│   │   │   ├── health.go           # Health, live, ready probes
│   │   │   ├── errors.go           # Shared service→API error mapping
│   │   │   └── validate.go         # Request body validation helper
│   │   └── middleware/
│   │       ├── request_id.go       # X-Request-ID generation/propagation
│   │       ├── timing.go           # X-Process-Time, slow request logging
│   │       ├── security_headers.go # HSTS, CSP, X-Frame-Options, etc.
│   │       ├── trusted_host.go     # Host header validation
│   │       ├── cors.go             # CORS configuration
│   │       ├── rate_limiter.go     # Redis + in-memory rate limiting
│   │       ├── error_handler.go    # Structured JSON error responses
│   │       └── auth.go             # JWT Bearer token validation
│   │
│   ├── config/
│   │   └── config.go               # Viper-based config with env vars
│   │
│   ├── db/
│   │   ├── postgres.go             # PostgreSQL connection with pool config
│   │   ├── redis.go                # Redis client with graceful fallback
│   │   ├── factory.go              # Database provider factory
│   │   ├── migrate.go              # Embedded migration runner
│   │   └── migrations/
│   │       ├── 000001_init.up.sql  # Initial schema
│   │       └── 000001_init.down.sql
│   │
│   ├── models/
│   │   └── user.go                 # GORM User model with soft delete
│   │
│   ├── dto/
│   │   ├── auth.go                 # Auth request/response DTOs
│   │   ├── users.go                # User CRUD DTOs
│   │   └── pagination.go           # Generic paginated response
│   │
│   ├── errors/
│   │   ├── api.go                  # HTTP error types with JSON envelope
│   │   ├── repository.go           # Data access error sentinels
│   │   └── service.go              # Business logic error sentinels
│   │
│   ├── repository/
│   │   ├── interfaces.go           # UserRepository interface
│   │   └── user_gorm.go            # GORM implementation (Postgres)
│   │
│   ├── service/
│   │   ├── auth.go                 # Auth logic, brute-force protection
│   │   └── users.go                # User CRUD, pagination, model↔DTO
│   │
│   ├── security/
│   │   ├── interfaces.go           # Hasher and TokenService interfaces
│   │   ├── hasher.go               # Bcrypt implementation
│   │   └── jwt.go                  # JWT token service
│   │
│   ├── logger/
│   │   └── logger.go               # zerolog setup, request-scoped logging
│   │
│   └── testutil/
│       └── testutil.go             # Test DB/Redis setup helpers
│
├── infra/
│   ├── Dockerfile                  # Multi-stage build
│   ├── docker-compose.local.yml    # Postgres, Redis, Nginx, API
│   ├── nginx.conf                  # Reverse proxy
│   ├── .env-example                # Environment variable template
│   └── commands/
│       ├── entrypoint.sh           # Container entrypoint
│       └── migrate.sh              # Migration runner
│
├── .github/
│   └── workflows/
│       └── ci.yml                  # GitHub Actions (lint, test)
│
├── .golangci.yml                   # Linter configuration
├── .gitignore
├── Makefile
├── go.mod
└── go.sum
```

## Getting Started

### Prerequisites

- **Go 1.22+**
- **PostgreSQL 17+**
- **Redis 7+**
- **Docker & Docker Compose** (optional, recommended)

### Option 1: Using Docker (Recommended)

1. **Clone the repository:**
```bash
git clone git@github.com:thealish/go-microservice.git
cd go-microservice
```

2. **Set up environment:**
```bash
cp infra/.env-example infra/.env
```

3. **Start all services:**
```bash
docker compose -f infra/docker-compose.local.yml up -d
```

This starts: API (port 8080), PostgreSQL, Redis, Nginx (port 80).

Migrations run automatically on startup.

4. **Verify:**
```bash
curl http://localhost/health
# {"status":"ok"}
```

### Option 2: Local Development

1. **Clone and configure:**
```bash
git clone git@github.com:thealish/go-microservice.git
cd go-microservice
cp infra/.env-example infra/.env
```

2. **Edit `infra/.env`** — set `POSTGRES_DSN` to your local Postgres:
```
POSTGRES_DSN=postgres://devuser:devpassword@localhost:5432/app_dev?sslmode=disable
REDIS_HOST=localhost
```

3. **Start backing services:**
```bash
docker run -d --name pg -p 5432:5432 \
  -e POSTGRES_USER=devuser \
  -e POSTGRES_PASSWORD=devpassword \
  -e POSTGRES_DB=app_dev \
  postgres:17

docker run -d --name redis -p 6379:6379 redis:7-alpine
```

4. **Run the application:**
```bash
go run ./cmd/server
```

5. **Verify:**
```bash
curl http://localhost:8080/
# {"app":"go-microservice","version":"0.1.0","status":"running"}
```

## Configuration

Configuration is managed through environment variables, loaded via Viper from `infra/.env`.

### Configuration Areas

| Section | Variables | Description |
|---------|-----------|-------------|
| Server | `SERVER_HOST`, `SERVER_PORT`, `SERVER_ENVIRONMENT` | Host, port, environment profile |
| JWT | `JWT_SECRET_KEY`, `JWT_ACCESS_TOKEN_EXPIRY` | Token signing and expiry |
| Auth | `AUTH_MAX_ATTEMPTS`, `AUTH_LOCKOUT_SECONDS` | Brute-force protection |
| Postgres | `POSTGRES_DSN`, `POSTGRES_POOL_SIZE` | Database connection and pooling |
| Redis | `REDIS_HOST`, `REDIS_PORT` | Cache and rate limiting backend |
| Rate Limit | `RATE_LIMIT_ENABLED`, `RATE_LIMIT_LIMIT_GET` | Per-IP request throttling |
| CORS | `CORS_ALLOWED_ORIGINS` | Cross-origin policy |
| Logging | `LOGGING_LEVEL`, `LOGGING_SLOW_REQUEST_THRESHOLD_MS` | Log verbosity and slow detection |

### Environment Profiles

- `development` — Pretty console logging, relaxed validation
- `test` — Test-friendly defaults
- `staging` / `production` — JSON logging, strict validation (no wildcard trusted hosts, JWT secret required)

## Database Migrations

Migrations use [golang-migrate](https://github.com/golang-migrate/migrate) with SQL files embedded in the binary.

Migrations run **automatically on startup** for PostgreSQL.

### Manual Migration Commands

```bash
# Apply migrations
make migrate-up

# Rollback last migration
make migrate-down

# Create new migration
make migrate-create name=add_user_profile
```

Migration files are stored in `internal/db/migrations/` and embedded via `//go:embed`.

## Development Commands

```bash
# Run & Build
make run                    # Start the application
make build                  # Build binary to bin/

# Testing
make test                   # Run all tests
make test-unit              # Run unit tests only
make test-integration       # Run integration tests

# Code Quality
make lint                   # Run golangci-lint
make fmt                    # Format code (gofmt + goimports)

# Database
make migrate-up             # Apply migrations
make migrate-down           # Rollback last migration
make migrate-create name=x  # Create new migration

# Docker
make docker-up              # Start PostgreSQL + Redis
```

## API Endpoints

### Health Probes

```
GET  /             # Service info (name, version, status)
GET  /health       # Health check
GET  /live         # Liveness probe
GET  /ready        # Readiness probe (DB + Redis)
```

### Authentication

```
POST /api/v1/auth/register   # Register new user
POST /api/v1/auth/login      # Login, get access + refresh tokens
POST /api/v1/auth/refresh    # Refresh access token
GET  /api/v1/auth/me         # Get current user (protected)
```

### Users (all protected)

```
GET    /api/v1/users           # List users (paginated: ?page=1&per_page=20)
GET    /api/v1/users/:id       # Get user by ID
POST   /api/v1/users           # Create user
PATCH  /api/v1/users/:id       # Update user (email/password)
DELETE /api/v1/users/:id       # Soft delete user
```

### Error Response Format

All errors return a structured JSON envelope:

```json
{
  "error": {
    "type": "UNAUTHORIZED",
    "message": "Invalid credentials",
    "request_id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

## Architecture

### Layered Architecture

```
Request → Middleware → Router → Handler → Service → Repository → Database
                                                       ↓
Response ← Middleware ← Router ← Handler ← DTO ← Model
```

1. **Handlers** (`internal/api/handlers/`) — Parse HTTP requests, call services, return responses
2. **Services** (`internal/service/`) — Business logic, error mapping, model↔DTO conversion
3. **Repositories** (`internal/repository/`) — Data access via GORM, error wrapping
4. **Models** (`internal/models/`) — GORM models with hooks and constraints
5. **DTOs** (`internal/dto/`) — Request/response types with validation tags

### Design Patterns

- **Repository Pattern** — Abstract data access behind `UserRepository` interface
- **Constructor Injection** — Dependencies wired explicitly in `cmd/server/main.go`
- **Interface Segregation** — Handlers define their own service interfaces (`AuthServicer`, `UserServicer`)
- **Sentinel Errors** — Typed errors at each layer for clean error propagation
- **Middleware Chain** — Composable request/response processing

### Middleware Stack (execution order)

1. Request ID → 2. Timing → 3. Security Headers → 4. Trusted Host → 5. CORS → 6. Rate Limiter → 7. Error Handler

## Security

### Authentication

- **JWT Tokens** — Access (15min) + Refresh (7 days) token pattern
- **Bcrypt** — Password hashing with `golang.org/x/crypto/bcrypt`
- **Brute-Force Protection** — Per-email + per-IP attempt tracking via Redis
- **Account Lockout** — Configurable max attempts, window, and lockout duration

### Security Headers

- `Strict-Transport-Security` — HSTS with includeSubDomains
- `X-Frame-Options: DENY` — Clickjacking protection
- `Content-Security-Policy: default-src 'self'` — XSS mitigation
- `X-Content-Type-Options: nosniff` — MIME sniffing protection
- `Referrer-Policy: strict-origin-when-cross-origin`
- `Server` and `X-Powered-By` headers stripped

### Additional Security

- **Input Validation** — Struct tag validation on all request DTOs
- **Rate Limiting** — Per-IP with separate GET/write limits
- **Trusted Hosts** — Host header whitelist (rejects unknown hosts)
- **Soft Deletion** — Data preserved, not permanently removed
- **SQL Injection Protection** — GORM parameterized queries
- **Graceful Redis Fallback** — Rate limiting continues with in-memory counters

## Testing

```bash
# Run all unit tests
make test

# Run with verbose output
go test ./... -v -count=1

# Run integration tests (requires Postgres + Redis)
go test -tags integration ./... -v -count=1

# Run specific package tests
go test ./internal/service/ -v
go test ./internal/security/ -v
go test ./internal/api/handlers/ -v
go test ./internal/api/middleware/ -v
```

### Test Organization

- **Unit tests** — Live next to source files (`*_test.go`), use testify mocks
- **Integration tests** — Behind `//go:build integration` build tag, require real Postgres + Redis
- **Test helpers** — `internal/testutil/` for shared DB/Redis setup

### Current Coverage

| Package | Tests |
|---------|-------|
| `internal/security` | 7 (hasher + JWT) |
| `internal/service` | 15 (auth + users) |
| `internal/api/handlers` | 13 (health + auth + users) |
| `internal/api/middleware` | 8 (rate limiter + trusted host + auth) |
| `internal/repository` | 5 (integration, Postgres) |

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Run tests: `make test`
5. Run linter: `make lint`
6. Format code: `make fmt`
7. Commit: `git commit -m 'Add amazing feature'`
8. Push: `git push origin feature/amazing-feature`
9. Open a Pull Request

### Code Style

- Follow Go conventions (camelCase unexported, PascalCase exported)
- Use `context.Context` for IO operations
- Early return error handling, no nesting
- No panics — return errors
- Define interfaces where consumed
- Keep files small and focused

## Acknowledgments

Built with:
- [Fiber](https://gofiber.io/) — HTTP framework
- [GORM](https://gorm.io/) — ORM
- [golang-migrate](https://github.com/golang-migrate/migrate) — Database migrations
- [zerolog](https://github.com/rs/zerolog) — Structured logging
- [go-redis](https://github.com/redis/go-redis) — Redis client
- [golang-jwt](https://github.com/golang-jwt/jwt) — JWT tokens
- [Viper](https://github.com/spf13/viper) — Configuration
- [testify](https://github.com/stretchr/testify) — Testing
- [validator](https://github.com/go-playground/validator) — Input validation
