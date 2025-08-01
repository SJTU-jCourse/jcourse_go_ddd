package apperror

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestErrorCreation(t *testing.T) {
	// Test predefined errors
	err := ErrNotFound
	if err.Code != 1001 {
		t.Errorf("Expected code 1001, got %d", err.Code)
	}
	if err.Category != CategoryDomain {
		t.Errorf("Expected category %d, got %d", CategoryDomain, err.Category)
	}

	// Test error creation with helpers
	customErr := NewValidationError(4005, "custom validation error")
	if customErr.Code != 4005 {
		t.Errorf("Expected code 4005, got %d", customErr.Code)
	}
	if customErr.Category != CategoryValidation {
		t.Errorf("Expected category %d, got %d", CategoryValidation, customErr.Category)
	}
}

func TestErrorWrapping(t *testing.T) {
	originalErr := errors.New("database connection failed")
	
	// Test WrapDB
	wrappedErr := WrapDB(originalErr)
	if !errors.Is(wrappedErr, originalErr) {
		t.Error("Wrapped error should be the same as original")
	}
	if wrappedErr.Category != CategoryInfrastructure {
		t.Error("Wrapped DB error should be infrastructure category")
	}

	// Test Wrap method
	baseErr := ErrNotFound
	wrappedWithMsg := baseErr.Wrap(originalErr)
	if wrappedWithMsg.Message != baseErr.Message {
		t.Error("Wrapped error should preserve original message")
	}
}

func TestErrorMetadata(t *testing.T) {
	err := ErrDB.
		WithMetadata("operation", "user_query").
		WithMetadata("user_id", 123).
		WithMetadata("timestamp", time.Now())

	if err.Metadata["operation"] != "user_query" {
		t.Error("Metadata should contain operation")
	}
	if err.Metadata["user_id"] != 123 {
		t.Error("Metadata should contain user_id")
	}
}

func TestErrorMessage(t *testing.T) {
	// Test error without wrapped error
	err := ErrNotFound
	if err.Error() != "resource not found" {
		t.Errorf("Expected 'resource not found', got '%s'", err.Error())
	}

	// Test error with wrapped error
	wrappedErr := ErrNotFound.Wrap(errors.New("user not found"))
	expected := "resource not found: user not found"
	if wrappedErr.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, wrappedErr.Error())
	}
}

func TestErrorMatching(t *testing.T) {
	// Test errors.Is
	err := ErrNotFound.WithMessage("custom message")
	if !errors.Is(err, ErrNotFound) {
		t.Error("Custom message error should still match ErrNotFound")
	}

	// Test errors.As
	var appErr *AppError
	if !errors.As(err, &appErr) {
		t.Error("Should be able to extract AppError")
	}
	if appErr.Message != "custom message" {
		t.Errorf("Expected 'custom message', got '%s'", appErr.Message)
	}
}

func TestUserMessage(t *testing.T) {
	// Test default user message
	err := ErrPermission
	if err.getUserMessage() != "You don't have permission to perform this action" {
		t.Error("Default permission user message should be user-friendly")
	}

	// Test custom user message
	customErr := ErrPermission.WithUserMessage("Custom permission message")
	if customErr.UserMessage != "Custom permission message" {
		t.Error("Custom user message should be preserved")
	}
}

func TestHTTPStatus(t *testing.T) {
	tests := []struct {
		err      *AppError
		expected int
	}{
		{ErrNotFound, 404},
		{ErrWrongInput, 400},
		{ErrWrongAuth, 401},
		{ErrPermission, 403},
		{ErrDB, 500},
	}

	for _, test := range tests {
		actual := test.err.HTTPStatus()
		if actual != test.expected {
			t.Errorf("Expected HTTP status %d for error %s, got %d", 
				test.expected, test.err.Message, actual)
		}
	}
}

func TestRetryable(t *testing.T) {
	// Test retryable errors
	retryableErr := ErrDB
	if !retryableErr.IsRetryable() {
		t.Error("DB errors should be retryable")
	}

	// Test non-retryable errors
	nonRetryableErr := ErrPermission
	if nonRetryableErr.IsRetryable() {
		t.Error("Permission errors should not be retryable")
	}
}

func TestSeverity(t *testing.T) {
	tests := []struct {
		err      *AppError
		expected string
	}{
		{ErrValidation, "low"},
		{ErrWrongAuth, "medium"},
		{ErrPermission, "medium"},
		{ErrNotFound, "high"},
		{ErrDB, "critical"},
	}

	for _, test := range tests {
		if test.err.GetSeverity() != test.expected {
			t.Errorf("Expected severity '%s' for error %s, got '%s'", 
				test.expected, test.err.Message, test.err.GetSeverity())
		}
	}
}

func TestErrorHandler(t *testing.T) {
	handler := NewErrorHandler(false, false)
	ctx := context.Background()

	// Test handling AppError
	appErr := ErrNotFound.WithMessage("user not found")
	response := handler.Handle(ctx, appErr)

	if response.Success {
		t.Error("Error response should not be successful")
	}
	if response.Code != ErrNotFound.Code {
		t.Errorf("Expected code %d, got %d", ErrNotFound.Code, response.Code)
	}
	if response.UserMessage == "" {
		t.Error("Response should have user message")
	}

	// Test handling generic error
	genericErr := errors.New("generic error")
	response = handler.Handle(ctx, genericErr)

	if response.Success {
		t.Error("Error response should not be successful")
	}
	if response.Code != 5002 { // WrapInternal creates code 5002
		t.Errorf("Expected code %d, got %d", 5002, response.Code)
	}
}

func TestErrorWithContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), "request_id", "req-123")
	ctx = context.WithValue(ctx, "user_id", 456)

	err := ErrDB.WithContext(ctx)

	if err.Metadata["request_id"] != "req-123" {
		t.Error("Error should contain request_id from context")
	}
	if err.Metadata["user_id"] != 456 {
		t.Error("Error should contain user_id from context")
	}
}

func TestErrorWithMessage(t *testing.T) {
	originalErr := ErrDB.WithMetadata("key", "value")
	newErr := originalErr.WithMessage("new message")

	if newErr.Message != "new message" {
		t.Errorf("Expected message 'new message', got '%s'", newErr.Message)
	}
	if newErr.Metadata["key"] != "value" {
		t.Error("Metadata should be preserved")
	}
	if newErr.Code != originalErr.Code {
		t.Error("Code should be preserved")
	}
	if newErr.Category != originalErr.Category {
		t.Error("Category should be preserved")
	}
}

func TestErrorWithUserMessage(t *testing.T) {
	originalErr := ErrDB
	newErr := originalErr.WithUserMessage("User friendly message")

	if newErr.UserMessage != "User friendly message" {
		t.Errorf("Expected user message 'User friendly message', got '%s'", newErr.UserMessage)
	}
	if newErr.Code != originalErr.Code {
		t.Error("Code should be preserved")
	}
}

func TestErrorChaining(t *testing.T) {
	// Test error chaining through multiple layers
	dbErr := errors.New("connection timeout")
	wrapped1 := WrapDB(dbErr).WithMetadata("layer", "repository")
	wrapped2 := wrapped1.WithMetadata("layer", "service").WithMessage("service error")
	wrapped3 := wrapped2.WithMetadata("layer", "handler").WithUserMessage("Service unavailable")

	if !errors.Is(wrapped3, dbErr) {
		t.Error("Final error should chain back to original DB error")
	}
	if wrapped3.Metadata["layer"] != "handler" {
		t.Error("Latest metadata should be preserved")
	}
	if wrapped3.UserMessage != "Service unavailable" {
		t.Error("User message should be preserved")
	}
}