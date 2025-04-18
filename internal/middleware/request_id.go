package middleware

import (
	"context"
	"net/http"
	"time"

	"drive/internal/util"

	"github.com/google/uuid"
)

const (
	requestIDKey contextKey = "requestID"
)

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware(logger *util.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Generate a new request ID
			requestID := uuid.New().String()

			// Add the request ID to the response header
			w.Header().Set("X-Request-ID", requestID)

			// Add the request ID to the request context
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)

			// Create a logger with request context
			requestLogger := logger.WithRequestID(requestID).
				WithMethod(r.Method).
				WithPath(r.URL.Path).
				WithRemoteAddr(r.RemoteAddr)

			// Log the request start
			requestLogger.Info("Request started")

			// Create a response writer that tracks the status code
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Start timing the request
			start := time.Now()

			// Call the next handler
			next.ServeHTTP(rw, r.WithContext(ctx))

			// Log the request completion
			duration := time.Since(start)
			requestLogger.WithDuration(duration).
				WithStatusCode(rw.statusCode).
				Info("Request completed")
		})
	}
}

// GetRequestID returns the request ID from the context
func GetRequestID(r *http.Request) string {
	if id, ok := r.Context().Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

// responseWriter is a wrapper around http.ResponseWriter that tracks the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
