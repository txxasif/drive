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
- OAuth authentication with Google and Facebook
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

### OAuth Authentication

- `POST /api/auth/oauth/login` - Authenticate with OAuth providers (Google, Facebook)

### Health Check

- `GET /health` - Service health check

## OAuth Integration

The application supports social login via Google and Facebook OAuth:

### Setup OAuth Providers

1. **Google OAuth Setup**:

   - Create a project in the [Google Developer Console](https://console.developers.google.com/)
   - Configure OAuth consent screen
   - Create OAuth client ID credentials for a web application
   - Add authorized redirect URIs for your application
   - Update the `.env` file with your Google credentials:
     ```
     GOOGLE_CLIENT_ID=your-client-id
     GOOGLE_CLIENT_SECRET=your-client-secret
     ```

2. **Facebook OAuth Setup**:
   - Create an app in the [Facebook Developer Portal](https://developers.facebook.com/)
   - Add Facebook Login product to your app
   - Configure Valid OAuth Redirect URIs
   - Update the `.env` file with your Facebook credentials:
     ```
     FACEBOOK_APP_ID=your-app-id
     FACEBOOK_APP_SECRET=your-app-secret
     ```

### Using OAuth in Your Application

The backend expects the client (frontend) to handle the initial OAuth flow:

1. Client initiates OAuth flow with the provider (Google/Facebook)
2. Client receives the access token from the provider
3. Client sends the token to the backend API endpoint (`/api/auth/oauth/login`)
4. Backend validates the token with the provider and:
   - Creates a new user account if the email doesn't exist
   - Returns JWT tokens for existing users

Example request to authenticate with Google:

```json
POST /api/auth/oauth/login
{
  "token": "google-oauth-access-token",
  "provider": "google"
}
```

Example response:

```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "username": "user1234",
    "first_name": "John",
    "last_name": "Doe",
    "provider": "google",
    "storage_used": 0,
    "storage_limit": 15000,
    "created_at": "2023-09-10T15:30:45Z",
    "updated_at": "2023-09-10T15:30:45Z"
  },
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

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
