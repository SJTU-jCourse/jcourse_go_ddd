package apperror

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"time"
)

// ErrorHandler provides centralized error handling utilities
type ErrorHandler struct {
	enableStackTrace bool
	enableLogging    bool
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(enableStackTrace, enableLogging bool) *ErrorHandler {
	return &ErrorHandler{
		enableStackTrace: enableStackTrace,
		enableLogging:    enableLogging,
	}
}

// Handle processes an error and returns a structured response
func (h *ErrorHandler) Handle(ctx context.Context, err error) *ErrorResponse {
	if err == nil {
		return &ErrorResponse{
			Success: true,
			Code:    0,
			Message: "success",
		}
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return h.handleAppError(ctx, appErr)
	}

	return h.handleGenericError(ctx, err)
}

// handleAppError processes an AppError
func (h *ErrorHandler) handleAppError(ctx context.Context, err *AppError) *ErrorResponse {
	response := &ErrorResponse{
		Success:     false,
		Code:        err.Code,
		Message:     err.Message,
		Category:    err.Category,
		UserMessage: err.getUserMessage(),
		Timestamp:   time.Now(),
	}

	if h.enableLogging {
		h.logError(ctx, err)
	}

	if h.enableStackTrace {
		response.StackTrace = h.getStackTrace()
	}

	if err.Metadata != nil {
		response.Metadata = err.Metadata
	}

	return response
}

// handleGenericError processes a generic error
func (h *ErrorHandler) handleGenericError(ctx context.Context, err error) *ErrorResponse {
	wrappedErr := WrapInternal(err)
	
	response := &ErrorResponse{
		Success:     false,
		Code:        wrappedErr.Code,
		Message:     wrappedErr.Message,
		Category:    wrappedErr.Category,
		UserMessage: "An unexpected error occurred",
		Timestamp:   time.Now(),
	}

	if h.enableLogging {
		h.logError(ctx, wrappedErr)
	}

	if h.enableStackTrace {
		response.StackTrace = h.getStackTrace()
	}

	return response
}

// logError logs the error with context
func (h *ErrorHandler) logError(ctx context.Context, err *AppError) {
	// Implementation would integrate with logging system
	// For now, just format the error
	logEntry := fmt.Sprintf(
		"[%s] Error %d: %s",
		time.Now().Format(time.RFC3339),
		err.Code,
		err.Error(),
	)
	
	if err.Err != nil {
		logEntry += fmt.Sprintf(" | Wrapped: %v", err.Err)
	}
	
	// In a real implementation, this would use a proper logger
	fmt.Printf("ERROR: %s\n", logEntry)
}

// getStackTrace captures the current stack trace
func (h *ErrorHandler) getStackTrace() []string {
	if !h.enableStackTrace {
		return nil
	}

	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	var stack []string
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		stack = append(stack, fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function))
	}

	return stack
}

// ErrorResponse represents a structured error response
type ErrorResponse struct {
	Success     bool            `json:"success"`
	Code        int             `json:"code"`
	Message     string          `json:"message"`
	Category    int             `json:"category,omitempty"`
	UserMessage string          `json:"user_message,omitempty"`
	Metadata    map[string]any  `json:"metadata,omitempty"`
	StackTrace  []string        `json:"stack_trace,omitempty"`
	Timestamp   time.Time       `json:"timestamp"`
}

// getUserMessage returns the user-friendly message or falls back to the default
func (e *AppError) getUserMessage() string {
	if e.UserMessage != "" {
		return e.UserMessage
	}
	
	// Default user messages based on category
	switch e.Category {
	case CategoryValidation:
		return "Please check your input and try again"
	case CategoryAuth:
		return "Authentication failed. Please check your credentials"
	case CategoryPermission:
		return "You don't have permission to perform this action"
	case CategoryDomain:
		return "Operation cannot be completed"
	case CategoryInfrastructure:
		return "Service temporarily unavailable. Please try again later"
	default:
		return "An error occurred"
	}
}

// Global error handler instance
var DefaultErrorHandler = NewErrorHandler(true, true)

// HandleError is a convenience function for error handling
func HandleError(ctx context.Context, err error) *ErrorResponse {
	return DefaultErrorHandler.Handle(ctx, err)
}

// IsRetryable checks if an error is retryable
func (e *AppError) IsRetryable() bool {
	switch e.Category {
	case CategoryInfrastructure:
		return e.Code == 5001 || e.Code == 5002 // DB or Network errors
	case CategoryDomain:
		return e.Code == 1003 // Rate limit
	default:
		return false
	}
}

// GetSeverity returns the severity level of the error
func (e *AppError) GetSeverity() string {
	switch e.Category {
	case CategoryValidation:
		return "low"
	case CategoryAuth, CategoryPermission:
		return "medium"
	case CategoryDomain:
		return "high"
	case CategoryInfrastructure:
		return "critical"
	default:
		return "medium"
	}
}

// WithContext adds context information to the error
func (e *AppError) WithContext(ctx context.Context) *AppError {
	if ctx == nil {
		return e
	}
	
	newErr := e.WithMessage(e.Message)
	
	// Add request ID if available
	if requestID := ctx.Value("request_id"); requestID != nil {
		newErr = newErr.WithMetadata("request_id", requestID)
	}
	
	// Add user ID if available
	if userID := ctx.Value("user_id"); userID != nil {
		newErr = newErr.WithMetadata("user_id", userID)
	}
	
	return newErr
}