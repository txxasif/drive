package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"drive/internal/model"
	"drive/internal/service"
)

// contextKey is a custom type for context keys
type contextKey string

const userKey contextKey = "user"

// Auth is a middleware that checks for a valid JWT token in the Authorization header
func Auth(authService service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Check if the Authorization header has the correct format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
				return
			}

			// Get the token
			token := parts[1]
			if token == "" {
				http.Error(w, "Token required", http.StatusUnauthorized)
				return
			}

			// Verify the token and get the user
			user, err := authService.GetUserByToken(r.Context(), token)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Add the user to the request context
			ctx := context.WithValue(r.Context(), userKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext retrieves the user ID from the request context
func GetUserIDFromContext(r *http.Request) (uint, error) {
	user := r.Context().Value(userKey)
	if user == nil {
		return 0, errors.New("user not found in context")
	}

	// Type assert the user to the correct type (model.User)
	u, ok := user.(*model.User)
	if !ok {
		return 0, errors.New("invalid user type in context")
	}

	return u.ID, nil
}
