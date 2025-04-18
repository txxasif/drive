package response

import (
	"encoding/json"
	"net/http"
)

// Standard response structure
type Response struct {
	Success bool           `json:"success"`
	Data    interface{}    `json:"data,omitempty"`
	Error   *ErrorResponse `json:"error,omitempty"`
	Meta    *MetaData      `json:"meta,omitempty"`
}

// Error response structure
type ErrorResponse struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details []string          `json:"details,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
}

// FieldError represents a validation error for a specific field
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// MetaData for pagination and other metadata
type MetaData struct {
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"per_page,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}

// Standard error codes
const (
	ErrBadRequest     = "BAD_REQUEST"
	ErrUnauthorized   = "UNAUTHORIZED"
	ErrForbidden      = "FORBIDDEN"
	ErrNotFound       = "NOT_FOUND"
	ErrInternalServer = "INTERNAL_SERVER_ERROR"
	ErrValidation     = "VALIDATION_ERROR"
	ErrDuplicateEntry = "DUPLICATE_ENTRY"
)

// Helper functions for common responses
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	response := Response{
		Success: statusCode >= 200 && statusCode < 300,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func Error(w http.ResponseWriter, statusCode int, code string, message string, details ...string) {
	response := Response{
		Success: false,
		Error: &ErrorResponse{
			Code:    code,
			Message: message,
			Details: details,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// ErrorWithFields sends an error response with field-specific errors
func ErrorWithFields(w http.ResponseWriter, statusCode int, code string, message string, fields map[string]string) {
	response := Response{
		Success: false,
		Error: &ErrorResponse{
			Code:    code,
			Message: message,
			Fields:  fields,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// Convenience methods for common error responses
func BadRequest(w http.ResponseWriter, message string, details ...string) {
	Error(w, http.StatusBadRequest, ErrBadRequest, message, details...)
}

func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, ErrUnauthorized, message)
}

func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, ErrForbidden, message)
}

func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, ErrNotFound, message)
}

func InternalError(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, ErrInternalServer, "An internal server error occurred")
}

func ValidationError(w http.ResponseWriter, details ...string) {
	Error(w, http.StatusBadRequest, ErrValidation, "Validation failed", details...)
}

// ValidationErrorWithFields sends a validation error response with field-specific errors
func ValidationErrorWithFields(w http.ResponseWriter, fields map[string]string) {
	ErrorWithFields(w, http.StatusBadRequest, ErrValidation, "Validation failed", fields)
}

// WithPagination adds pagination metadata to the response
func WithPagination(w http.ResponseWriter, statusCode int, data interface{}, page, perPage, totalCount int) {
	totalPages := (totalCount + perPage - 1) / perPage

	response := Response{
		Success: true,
		Data:    data,
		Meta: &MetaData{
			Page:       page,
			PerPage:    perPage,
			TotalPages: totalPages,
			TotalCount: totalCount,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
