package apperror

import (
	"fmt"
	"net/http"
)

// Error categories
const (
	CategoryDomain         = 1000 // Business logic errors
	CategoryAuth           = 2000 // Authentication errors
	CategoryPermission     = 3000 // Authorization errors
	CategoryValidation     = 4000 // Input validation errors
	CategoryInfrastructure = 5000 // Infrastructure errors
)

// AppError represents a structured application error
type AppError struct {
	Code        int            // Error code
	Message     string         // Error message
	Category    int            // Error category
	Err         error          // Wrapped error for context
	Metadata    map[string]any // Additional error context
	httpStatus  int            // HTTP status code (internal use)
	UserMessage string         // User-friendly message
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// Is checks if the error matches the target error
func (e *AppError) Is(target error) bool {
	if other, ok := target.(*AppError); ok {
		return e.Code == other.Code
	}
	return false
}

// Wrap wraps an error with additional context
func (e *AppError) Wrap(err error) *AppError {
	return &AppError{
		Code:        e.Code,
		Message:     e.Message,
		Category:    e.Category,
		Err:         err,
		Metadata:    e.Metadata,
		httpStatus:  e.httpStatus,
		UserMessage: e.UserMessage,
	}
}

// WithMessage creates a new error with a custom message
func (e *AppError) WithMessage(message string) *AppError {
	return &AppError{
		Code:        e.Code,
		Message:     message,
		Category:    e.Category,
		Err:         e.Err,
		Metadata:    e.Metadata,
		httpStatus:  e.httpStatus,
		UserMessage: e.UserMessage,
	}
}

// WithMetadata adds metadata to the error
func (e *AppError) WithMetadata(key string, value any) *AppError {
	newErr := e.WithMessage(e.Message)
	if newErr.Metadata == nil {
		newErr.Metadata = make(map[string]any)
	}
	newErr.Metadata[key] = value
	return newErr
}

// WithUserMessage sets a user-friendly message
func (e *AppError) WithUserMessage(message string) *AppError {
	newErr := e.WithMessage(e.Message)
	newErr.UserMessage = message
	return newErr
}

// HTTPStatus returns the appropriate HTTP status code
func (e *AppError) HTTPStatus() int {
	if e.httpStatus != 0 {
		return e.httpStatus
	}

	// Default mapping based on category
	switch e.Category {
	case CategoryDomain:
		switch e.Code {
		case 1001: // ErrNotFound
			return http.StatusNotFound
		case 1002: // ErrWrongInput
			return http.StatusBadRequest
		case 1003: // ErrRateLimit
			return http.StatusTooManyRequests
		default:
			return http.StatusBadRequest
		}
	case CategoryAuth:
		switch e.Code {
		case 2002, 2003: // ErrWrongAuth, ErrSession
			return http.StatusUnauthorized
		case 2004: // ErrSuspended
			return http.StatusForbidden
		default:
			return http.StatusUnauthorized
		}
	case CategoryPermission:
		return http.StatusForbidden
	case CategoryValidation:
		return http.StatusBadRequest
	case CategoryInfrastructure:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// Helper functions for creating structured errors

// NewDomainError creates a domain error
func NewDomainError(code int, message string) *AppError {
	return &AppError{
		Code:     code,
		Message:  message,
		Category: CategoryDomain,
	}
}

// NewAuthError creates an authentication error
func NewAuthError(code int, message string) *AppError {
	return &AppError{
		Code:     code,
		Message:  message,
		Category: CategoryAuth,
	}
}

// NewPermissionError creates a permission error
func NewPermissionError(code int, message string) *AppError {
	return &AppError{
		Code:     code,
		Message:  message,
		Category: CategoryPermission,
	}
}

// NewValidationError creates a validation error
func NewValidationError(code int, message string) *AppError {
	return &AppError{
		Code:     code,
		Message:  message,
		Category: CategoryValidation,
	}
}

// NewInfrastructureError creates an infrastructure error
func NewInfrastructureError(code int, message string) *AppError {
	return &AppError{
		Code:     code,
		Message:  message,
		Category: CategoryInfrastructure,
	}
}

// WrapDB wraps a database error
func WrapDB(err error) *AppError {
	return NewInfrastructureError(5001, "database error").Wrap(err)
}

// WrapInternal wraps an internal error
func WrapInternal(err error) *AppError {
	return NewInfrastructureError(5002, "internal error").Wrap(err)
}

// Predefined error constants
var (
	// Domain Errors (1000-1999)
	ErrNotFound     = NewDomainError(1001, "resource not found")
	ErrWrongInput   = NewDomainError(1002, "invalid input")
	ErrRateLimit    = NewDomainError(1003, "rate limit exceeded")
	ErrInvalidParam = NewDomainError(1004, "invalid parameter")

	// Authentication Errors (2000-2999)
	ErrExpired        = NewAuthError(2001, "token expired")
	ErrWrongAuth      = NewAuthError(2002, "invalid credentials")
	ErrSession        = NewAuthError(2003, "session error")
	ErrSuspended      = NewAuthError(2004, "account suspended")
	ErrNoTargetCourse = NewAuthError(2005, "target course not found")
	ErrNoSemester     = NewAuthError(2006, "semester not found")
	ErrUserNotFound   = NewAuthError(2007, "user not found")
	ErrInvalidToken   = NewAuthError(2008, "invalid token")

	// Permission Errors (3000-3999)
	ErrPermission   = NewPermissionError(3001, "permission denied")
	ErrAccessDenied = NewPermissionError(3002, "access denied")
	ErrUnauthorized = NewPermissionError(3003, "unauthorized access")

	// Validation Errors (4000-4999)
	ErrValidation      = NewValidationError(4001, "validation failed")
	ErrInvalidEmail    = NewValidationError(4002, "invalid email format")
	ErrInvalidPassword = NewValidationError(4003, "invalid password format")
	ErrRequiredField   = NewValidationError(4004, "required field missing")

	// Infrastructure Errors (5000-5999)
	ErrDB       = NewInfrastructureError(5001, "database error")
	ErrNetwork  = NewInfrastructureError(5002, "network error")
	ErrExternal = NewInfrastructureError(5003, "external service error")
	ErrInternal = NewInfrastructureError(5004, "internal server error")
)
