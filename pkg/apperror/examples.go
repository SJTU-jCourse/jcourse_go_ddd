package apperror

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

// Example usage and best practices for the error system

// Example 1: Basic error handling
func ExampleBasicErrorHandling() {
	// Using predefined errors
	err := processUserInput("")
	if err != nil {
		if errors.Is(err, ErrWrongInput) {
			// Handle specific error
			fmt.Println("Invalid input:", err)
		} else {
			// Handle other errors
			fmt.Println("Error:", err)
		}
	}
}

func processUserInput(input string) error {
	if input == "" {
		return ErrWrongInput.WithMessage("input cannot be empty")
	}
	return nil
}

// Example 2: Error wrapping with context
func ExampleErrorWrapping() {
	err := databaseOperation()
	if err != nil {
		// Wrap with additional context
		wrappedErr := ErrDB.Wrap(err).WithMetadata("operation", "user_query")
		fmt.Println(wrappedErr)
	}
}

func databaseOperation() error {
	return errors.New("connection timeout")
}

// Example 3: Custom error creation
func ExampleCustomError() {
	err := NewValidationError(4005, "email already exists").
		WithUserMessage("This email address is already registered").
		WithMetadata("field", "email").
		WithMetadata("email", "user@example.com")

	response := HandleError(context.Background(), err)
	fmt.Printf("Error response: %+v\n", response)
}

// Example 4: Error handling in HTTP handlers
func ExampleHTTPErrorHandling() {
	ctx := context.WithValue(context.Background(), "request_id", "req-123")

	err := performRestrictedOperation(ctx)
	response := HandleError(ctx, err)

	fmt.Printf("HTTP Response: Code=%d, Message=%s, UserMessage=%s\n",
		response.Code, response.Message, response.UserMessage)
}

func performRestrictedOperation(ctx context.Context) error {
	// Simulate permission error
	return ErrPermission.WithMessage("user cannot delete this resource").
		WithContext(ctx)
}

// Example 5: Error handling with middleware
func ExampleMiddlewareErrorHandling() {
	err := processData("invalid-data")

	// Check error type and handle accordingly
	var appErr *AppError
	if errors.As(err, &appErr) {
		switch appErr.Category {
		case CategoryValidation:
			fmt.Println("Validation error:", appErr.UserMessage)
		case CategoryAuth:
			fmt.Println("Auth error:", appErr.UserMessage)
		case CategoryPermission:
			fmt.Println("Permission error:", appErr.UserMessage)
		default:
			fmt.Println("Other error:", appErr.Message)
		}
	}
}

func processData(data string) error {
	if data == "invalid-data" {
		return NewValidationError(4006, "invalid data format").
			WithUserMessage("The data format is incorrect")
	}
	return nil
}

// Example 6: Error handling in tests
func TestErrorHandling(t *testing.T) {
	err := someOperation()

	// Test for specific error type
	if errors.Is(err, ErrNotFound) {
		t.Log("Correctly identified not found error")
		return
	}

	// Test for error category
	var appErr *AppError
	if errors.As(err, &appErr) && appErr.Category == CategoryValidation {
		t.Log("Correctly identified validation error")
		return
	}

	t.Errorf("Unexpected error: %v", err)
}

func someOperation() error {
	return ErrNotFound.WithMessage("user not found")
}

// Example 7: Error propagation
func ExampleErrorPropagation() {
	err := serviceLayer()
	if err != nil {
		// The error will have all the context from each layer
		response := HandleError(context.Background(), err)
		fmt.Printf("Final error response: %+v\n", response)
	}
}

func serviceLayer() error {
	err := repositoryLayer()
	if err != nil {
		return ErrDB.Wrap(err).WithMetadata("service", "user_service")
	}
	return nil
}

func repositoryLayer() error {
	return errors.New("database connection failed")
}

// Example 8: Error recovery and retry logic
func ExampleRetryLogic() {
	err := retryOperation()
	if err != nil {
		var appErr *AppError
		if errors.As(err, &appErr) && appErr.IsRetryable() {
			fmt.Println("Error is retryable:", appErr.Message)
		} else {
			fmt.Println("Error is not retryable:", appErr.Message)
		}
	}
}

func retryOperation() error {
	// Simulate a retryable error
	return ErrDB.WithMessage("temporary database failure")
}

// Example 9: Error logging and monitoring
func ExampleErrorLogging() {
	ctx := context.Background()
	err := criticalOperation()

	response := HandleError(ctx, err)

	// Log the error with severity
	var appErr *AppError
	if errors.As(err, &appErr) {
		fmt.Printf("Error severity: %s\n", appErr.GetSeverity())
		fmt.Printf("Error category: %d\n", appErr.Category)
		fmt.Printf("Error metadata: %+v\n", response.Metadata)
	}
}

func criticalOperation() error {
	return ErrInternal.WithMessage("critical system failure").
		WithMetadata("component", "payment_service").
		WithMetadata("severity", "critical")
}

// Example 10: Error handling patterns
func ExampleErrorPatterns() {
	// Pattern 1: Early return with error
	if err := validateInput("test"); err != nil {
		return // Handle error
	}

	// Pattern 2: Error wrapping with context
	err := processRequest("test")
	if err != nil {
		// Wrap with additional context
		wrappedErr := ErrWrongInput.Wrap(err)
		fmt.Println(wrappedErr)
	}

	// Pattern 3: Error transformation
	err = processDataLayer("test")
	if err != nil {
		// Transform infrastructure error to user-friendly message
		userErr := ErrDB.WithUserMessage("Service temporarily unavailable")
		fmt.Println(userErr.UserMessage)
	}
}

func validateInput(input string) error {
	if input == "" {
		return ErrRequiredField.WithMessage("input is required")
	}
	return nil
}

func processRequest(input string) error {
	if input == "invalid" {
		return errors.New("invalid request format")
	}
	return nil
}

func processDataLayer(input string) error {
	return errors.New("database timeout")
}
