package errs

import (
	"fmt"
	"log/slog"
)

// AppError represents a structured application error
type AppError struct {
	Err      error             `json:"-"`               // Underlying error
	Code     string            `json:"code"`            // Machine readable code (e.g. "VALIDATION_FAILED", "UNAUTHORIZED")
	Msg      string            `json:"message"`         // User-friendly message
	Status   int               `json:"status"`          // HTTP status code
	Field    map[string]string `json:"field,omitempty"` // Field-specific errors
	LogLevel slog.Level        `json:"-"`               // Log level for the error
}

// Implements the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Msg)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// Helper to easily change severity on the fly
func (e *AppError) WithLevel(level slog.Level) *AppError {
	e.LogLevel = level
	return e
}

// New creates a new AppError
func New(err error, code, msg string, status int, field map[string]string) *AppError {
	return &AppError{
		Err:      err,
		Code:     code,
		Msg:      msg,
		Status:   status,
		Field:    field,
		LogLevel: slog.LevelInfo,
	}
}

// NewInternal creates a generic internal server error
func NewInternal(err error) *AppError {
	return &AppError{
		Err:      err,
		Code:     "INTERNAL_SERVER_ERROR",
		Msg:      "An unexpected error occurred. Please try again later.",
		Status:   500,
		LogLevel: slog.LevelError,
	}
}

// NewValidation creates a generic validation error
func NewValidation(err error) *AppError {
	return &AppError{
		Err:      err,
		Code:     "VALIDATION_FAILED",
		Msg:      "Please check the input fields for errors.",
		Status:   422,
		LogLevel: slog.LevelInfo,
	}
}

// NewUnauthorized creates an unauthorized error
func NewUnauthorized(err error) *AppError {
	return &AppError{
		Err:      err,
		Code:     "UNAUTHORIZED",
		Msg:      "You are not authorized to perform this action.",
		Status:   403,
		LogLevel: slog.LevelWarn,
	}
}

// NewNotFound creates a not found error
func NewNotFound(err error) *AppError {
	return &AppError{
		Err:      err,
		Code:     "NOT_FOUND",
		Msg:      "The requested resource was not found.",
		Status:   404,
		LogLevel: slog.LevelInfo,
	}
}

// NewBadRequest creates a bad request error
func NewBadRequest(err error) *AppError {
	return &AppError{
		Err:      err,
		Code:     "BAD_REQUEST",
		Msg:      "The request was invalid or cannot be served.",
		Status:   400,
		LogLevel: slog.LevelInfo,
	}
}
