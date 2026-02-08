# ElGopher

A modern, lightweight blog platform built with Go, featuring markdown support, user authentication, internationalization, and server-side rendering with HTMX.

## Overview

ElGopher is a full-featured blogging platform that emphasizes simplicity, performance, and developer experience. It provides a clean interface for creating and managing articles with markdown support, syntax highlighting for code blocks, and a responsive design that works across all devices.

## Features

- **Article Management**: Create, edit, publish, and manage blog articles with markdown syntax
- **Markdown Support**: Full markdown parsing with syntax highlighting for code blocks (using Goldmark and Chroma)
- **User Authentication**: Secure user login and session management with bcrypt password hashing
- **Internationalization (i18n)**: Multi-language support (English and Spanish)
- **Theme Support**: Light and dark theme preferences
- **HTMX Integration**: Dynamic, responsive UI without heavy JavaScript frameworks
- **Session Management**: PostgreSQL-backed session storage with automatic cleanup
- **Security Features**:
  - CSRF protection with nosurf
  - Secure session cookies
  - Content Security Policy headers
  - XSS protection
  - HTTP Strict Transport Security (HSTS) in production
- **Database Migrations**: Structured database schema with versioned migrations
- **Docker Support**: Ready-to-deploy containerized application
- **Production Ready**: Configured for deployment on Fly.io

## Technology Stack

### Backend
- **Go 1.25**: Modern Go with latest features
- **pgx/v5**: High-performance PostgreSQL driver and connection pooling
- **alice**: Middleware chaining
- **scs**: Session management with PostgreSQL store
- **bcrypt**: Password hashing
- **goldmark**: Markdown parsing
- **chroma**: Syntax highlighting

### Frontend
- **HTMX**: Dynamic HTML without heavy JavaScript
- **Server-side templates**: Go's html/template engine
- **Responsive design**: Mobile-first approach

### Database
- **PostgreSQL**: Primary data store with connection pooling

## Project Structure

```
ElGopher/
├── cmd/
│   ├── web/           # Main web application
│   │   ├── main.go           # Application entry point
│   │   ├── handlers.go       # HTTP request handlers
│   │   ├── routes.go         # Route definitions
│   │   ├── middleware.go     # HTTP middleware
│   │   ├── templates.go      # Template rendering
│   │   ├── helpers.go        # Helper functions
│   │   └── context.go        # Request context management
│   └── seed/          # Database seeding utility
├── internal/
│   ├── models/        # Data models and database operations
│   │   ├── articles.go       # Article model
│   │   ├── users.go          # User model
│   │   └── mocks/            # Mock implementations for testing
│   ├── validator/     # Input validation utilities
│   ├── i18n/          # Internationalization support
│   └── assert/        # Testing assertions
├── database/
│   └── migrations/    # SQL migration files
├── ui/
│   ├── html/          # HTML templates
│   │   ├── base.html         # Base template
│   │   ├── pages/            # Page templates
│   │   └── partials/         # Reusable template partials
│   ├── static/        # Static assets (CSS, JS, images)
│   └── efs.go         # Embedded file system
├── Dockerfile         # Multi-stage Docker build
├── fly.toml          # Fly.io deployment configuration
└── go.mod            # Go module dependencies
```

## Database Schema

### Articles Table
- `id`: Primary key
- `title`: Article title
- `body`: Markdown content
- `slug`: URL-friendly identifier
- `excerpt`: Auto-generated article preview
- `is_published`: Publication status
- `created`: Creation timestamp
- `updated_at`: Last update timestamp

### Users Table
- `id`: Primary key
- `name`: User display name
- `email`: Unique user email
- `hashed_password`: Bcrypt-hashed password
- `created`: Registration timestamp

### Sessions Table
- PostgreSQL-backed session storage
- Automatic expiration (12-hour lifetime)

## Installation

### Prerequisites
- Go 1.25 or higher
- PostgreSQL 12 or higher
- Docker (optional, for containerized deployment)

### Local Development Setup

1. Clone the repository:
```bash
git clone https://github.com/AgustinPagotto/ElGopher.git
cd ElGopher
```

2. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Configure your database:
```bash
export DATABASE_URL="postgres://user:password@localhost:5432/elgopher"
```

4. Run database migrations:
```bash
# Apply migrations (using your preferred migration tool)
```

5. Install dependencies:
```bash
go mod download
```

6. Run the application:
```bash
go run ./cmd/web
```

The application will start on `http://localhost:4000` (or the port specified in your environment).

## Configuration

Environment variables:

- `DATABASE_URL`: PostgreSQL connection string (default: `postgres://postgres:admin@localhost:5432/elgopher`)
- `PORT`: Server port (default: `4000`)
- `ADMIN_EMAIL`: Initial admin user email
- `ADMIN_NAME`: Initial admin user name
- `ADMIN_PASSWORD`: Initial admin user password
- `IS_PROD`: Production mode flag

## Docker Deployment

Build and run with Docker:

```bash
docker build -t elgopher .
docker run -p 8080:8080 \
  -e DATABASE_URL="your_database_url" \
  -e IS_PROD="true" \
  elgopher
```

The Docker image uses a multi-stage build for optimal size and security, running as a non-root user.

## Fly.io Deployment

The application is configured for deployment on Fly.io:

```bash
fly launch
fly deploy
```

Configuration is provided in `fly.toml` with:
- Auto-start/stop machines
- 256MB memory allocation
- HTTPS enforcement
- São Paulo (GRU) region

## API Routes

### Public Routes
- `GET /`: Homepage
- `GET /about`: About page
- `GET /articles`: List all published articles
- `GET /article/view/{slug}`: View specific article
- `GET /projects`: Projects page
- `GET /user/login`: Login page
- `POST /user/login`: Login submission
- `GET /ping`: Health check endpoint

### Protected Routes (Authentication Required)
- `GET /article/create`: Create article form
- `POST /article/create`: Create article submission
- `GET /article/edit/{slug}`: Edit article form
- `PATCH /article/{id}`: Update article
- `POST /user/logout`: Logout

### Preference Routes
- `GET /pref/lng`: Set language preference
- `GET /pref/thm`: Set theme preference

## Testing

The project includes comprehensive test coverage:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./cmd/web
go test ./internal/models
```

Tests include:
- Handler tests with mock implementations
- Middleware tests
- Template rendering tests
- Model tests with test helpers

## Development Features

- **Hot Reload**: Use tools like `air` or `reflex` for automatic reloading during development
- **Mock Implementations**: Built-in mocks for testing without database dependencies
- **Logging**: Structured JSON logging with `slog`
- **Connection Pooling**: Efficient database connection management with pgxpool
- **Request Timeouts**: 5-second timeout middleware to prevent hanging requests
- **Panic Recovery**: Graceful error handling and recovery

## Security Considerations

- Passwords are hashed using bcrypt with cost factor 12
- CSRF protection on all state-changing requests
- Secure session cookies (Secure flag in production)
- Content Security Policy headers
- XSS protection headers
- Clickjacking prevention (X-Frame-Options)
- MIME type sniffing prevention
- HSTS in production environments

## Performance

- Lightweight binary (multi-stage Docker build with stripped symbols)
- Connection pooling for database efficiency
- Template caching for faster rendering
- Embedded file system for static assets
- Minimal external dependencies

## Contributing

Contributions are welcome! Please follow these guidelines:
1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Submit a pull request

## License

MIT License - Copyright (c) 2025 Agustin Oliveros Pagotto

See [LICENSE](LICENSE) file for full details.

## Author

**Agustin Oliveros Pagotto**

## Acknowledgments

Built with modern Go best practices and leveraging excellent open-source libraries:
- Goldmark for markdown parsing
- Chroma for syntax highlighting
- pgx for PostgreSQL connectivity
- SCS for session management
- Alice for middleware chaining