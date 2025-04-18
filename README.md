# Drive Storage System - User Service

This project implements a comprehensive user management system for the Drive Storage application, following SOLID principles and clean architecture.

## Architecture

The application follows a layered architecture pattern:

- **Model Layer**: Data structures representing database entities
- **Repository Layer**: Database operations and abstractions
- **Service Layer**: Business logic and use cases
- **Handler Layer**: HTTP request handling
- **Middleware**: Cross-cutting concerns like authentication
- **Routes**: API endpoint definitions

## Directory Structure

```
drive/
├── api/
│   └── routes/           # Route definitions
├── cmd/
│   └── server/           # Application entrypoint
├── internal/
│   ├── config/           # Configuration management
│   ├── database/         # Database connection utilities
│   ├── domain/           # Business domain interfaces and DTOs
│   ├── handler/          # HTTP request handlers
│   ├── middleware/       # HTTP middleware components
│   ├── model/            # Data models
│   ├── repository/       # Data access implementations
│   ├── service/          # Business logic implementations
│   └── util/             # Utility functions
├── .env                  # Environment variables
├── go.mod                # Go module definition
├── go.sum                # Go module checksum
└── README.md             # Project documentation
```

## Features

- User registration with email and password
- User authentication with JWT
- User profile management (get, update, delete)
- Secure password handling with bcrypt
- Token-based authentication
- Database integration with GORM
- Environment-based configuration
- Clean code separation following SOLID principles

## API Endpoints

### User Management

- `POST /api/users/register` - Register a new user
- `POST /api/users/login` - Login and get JWT token
- `GET /api/users/{id}` - Get user profile (requires authentication)
- `PUT /api/users/{id}` - Update user profile (requires authentication)
- `DELETE /api/users/{id}` - Delete user (requires authentication)

### Health Check

- `GET /health` - Service health check

## Getting Started

### Prerequisites

- Go 1.22 or later
- PostgreSQL

### Setup

1. Clone the repository
2. Create a PostgreSQL database
3. Copy `.env.example` to `.env` and configure environment variables
4. Run the application:

```bash
go run cmd/server/main.go
```

## SOLID Principles Implementation

- **Single Responsibility Principle**: Each component has a single responsibility (e.g., repositories for data access, services for business logic)
- **Open/Closed Principle**: Extension is possible without modifying existing code through interfaces
- **Liskov Substitution Principle**: Interface implementations can be substituted without affecting functionality
- **Interface Segregation Principle**: Small, focused interfaces for specific concerns
- **Dependency Inversion Principle**: High-level modules depend on abstractions, not concrete implementations

## Error Handling

The application implements consistent error handling patterns:

- Domain-specific errors at the service layer
- HTTP status code mapping in handlers
- Detailed error messages for debugging
- User-friendly error responses for clients

## Security Considerations

- Passwords are hashed using bcrypt
- JWT tokens with expiration
- Bearer token authentication
- HTTPS recommended for production
- Input validation for all requests

## Future Improvements

- Email verification
- Password reset functionality
- Role-based authorization
- Rate limiting
- API documentation with Swagger
- Integration with object storage for file management 