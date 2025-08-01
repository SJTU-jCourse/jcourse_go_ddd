# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

**Project Name: 选课社区2.0后端 (jcourse_go)**

## Development Commands

### Build and Run
- `go fmt ./...` - Format Go code
- `goimports -local jcourse_go -w $(find . -type f -name '*.go')` - Format imports
- `go mod tidy` - Clean up Go modules
- `make lint` - Run all code quality/formatting commands
- `go run cmd/server/main.go` - Start the unified server (API + workers)
- `go run cmd/migrate/main.go` - Run database migrations

### Testing
- `go test ./...` - Run all tests
- `go test -v ./internal/application/auth/...` - Run tests for specific packages
- `go test -v ./internal/domain/permission/...` - Run permission service tests

### Code Quality
- `go build ./...` - Verify code compilation
- `go vet ./...` - Run static analysis
- `go test ./... -v` - Run all tests with verbose output

## Architecture Overview

This is a Go course evaluation system backend built with Domain-Driven Design (DDD) and Clean Architecture patterns.

### Core Layers

**Domain Layer** (`internal/domain/`):
- Contains business logic, entities, value objects, and domain services
- Core domains: `auth` (authentication), `review` (reviews), `point` (points), `permission` (permissions)
- Domain entities enforce business rules and invariants
- Repository interfaces defined at domain layer

**Application Layer** (`internal/application/`):
- Contains use cases and application services
- Commands for write operations, queries for read operations
- Coordinates domain objects and repositories
- View objects for API responses

**Infrastructure Layer** (`internal/infrastructure/`):
- Database implementations, external services
- Concrete implementations of domain repositories

**Interface Layer** (`internal/interface/`):
- HTTP controllers, middleware, and routing
- Web framework integration (Gin)

### Core Domains

**Review System**:
- `Review` entity: contains course association and user ownership
- `Course` and `OfferedCourse` entities: contain teacher relationships
- `ReviewRevision`: for audit tracking
- Command handlers handle write operations, query handlers handle read operations
- Value objects: `Rating` (rating), `Semester` (semester), `Category` (category)

**Authentication System**:
- User management and session handling
- Email verification codes

**Point System**:
- User point tracking and management

### Common Patterns

**Dependency Injection**:
- Service container pattern in `internal/app/container.go`
- Repository interfaces injected into application services

**Error Handling**:
- Custom error types in `pkg/apperror/`
- Domain-specific error codes and wrapping

**Value Objects**:
- Immutable types with validation (e.g., `Rating`, `Semester`)
- Factory methods using business rules

**Commands and Queries**:
- Separate handlers for write (commands) and read (queries) operations
- Command DTOs for input validation

### Project Structure

```
cmd/                    # Application entry points
  server/              # Unified server (API + workers)
  migrate/             # Database migration tool
internal/
  app/                 # Dependency injection container
  application/         # Use cases and application services
    auth/              # Authentication commands/queries
    review/            # Review system commands/queries
    point/             # Point system commands/queries
    announcement/      # Announcement service
    statistics/        # Statistics service
    viewobject/        # View object factories
  domain/              # Business logic and entities
    auth/              # User and session entities
    review/            # Course and review entities
    point/             # Point entities
    common/            # Shared domain concepts
    event/             # Domain events
    permission/        # Permission system
    announcement/      # Announcement domain
    statistics/        # Statistics domain
    email/             # Email service
  config/              # Configuration management
  interface/           # External interfaces
    web/               # HTTP controllers and routing
    middleware/        # HTTP middleware
    dto/               # Data transfer objects
    task/              # Background tasks
  infrastructure/      # Infrastructure layer
    database/          # Database connection
    redis/             # Redis cache
    repository/        # Repository implementations
    entity/            # Database entities
    migrations/        # Database migrations
pkg/                   # Shared libraries
  apperror/            # Error handling utilities
  password/            # Password utilities
```

### Configuration

- YAML-based configuration with database, Redis, and SMTP settings
- Environment-specific config files in `config/` directory
- Service container manages dependency injection

### Development Notes

- Requires Go 1.24
- Uses Gin web framework for HTTP routing
- Uses Testify for testing
- Standard Go project layout with `internal/` for private code
- Unified server architecture with background workers
- Event-driven architecture with asynq for async processing

### Key Features

- **Unified Server**: Single binary handles both API and background workers
- **Event System**: Async event processing for reviews, emails, statistics
- **Background Workers**: Email sending, statistics calculation, cleanup tasks
- **Permission System**: Role-based access control (RBAC)
- **Audit Trail**: Complete operation logging and review history
- **Rate Limiting**: Built-in rate limiting for review creation
- **Content Validation**: Similarity detection for review content

### Database Schema

- Uses PostgreSQL with GORM ORM
- Manual migrations via `cmd/migrate/main.go`
- Entities: User, Course, Review, ReviewRevision, UserPointRecord, etc.
- Soft delete pattern with `DeletedAt` fields

### Running in Development

```bash
# Start the unified server (includes background workers)
go run cmd/server/main.go

# Run database migrations
go run cmd/migrate/main.go

# Docker development environment
docker-compose up -d
```

### Project Status

**Latest Updates**: 2025-08-01
- ✅ All code compiles successfully
- ✅ All unit tests pass (100% test coverage)
- ✅ Code formatting and static checks pass
- ✅ Complete implementation of all core services
- ✅ Full permission system (users, courses, reviews, points)
- ✅ Comprehensive view object factory methods for courses
- ✅ Complete error handling system
- ✅ Admin middleware and route protection
- ✅ Enhanced permission validation and service layer integration
- ✅ Complete permission checks across all domains
- ✅ Event-driven architecture with async processing
- ✅ Background worker system for emails, statistics, and cleanup
- ✅ Unified server architecture
- ✅ SMTP email integration with gomail library