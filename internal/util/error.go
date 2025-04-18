package util

import "fmt"

// AppError represents an application error
type AppError struct {
	Code    string
	Message string
	Details []string
}

// Error implements the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewAppError creates a new application error
func NewAppError(code, message string, details ...string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}
