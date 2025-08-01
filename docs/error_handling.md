# Error Handling System

## Overview

This document describes the systematic error handling approach used throughout the jcourse_go application. The system provides consistent error creation, handling, and response patterns.

## Error Categories

Errors are categorized into five main types:

1. **Domain Errors (1000-1999)**: Business logic errors
2. **Auth Errors (2000-2999)**: Authentication and authorization errors  
3. **Permission Errors (3000-3999)**: Permission and access control errors
4. **Validation Errors (4000-4999)**: Input validation errors
5. **Infrastructure Errors (5000-5999)**: System and infrastructure errors

## Error Structure

```go
type AppError struct {
    Code        int            // Error code
    Message     string         // Error message
    Category    int            // Error category
    Err         error          // Wrapped error for context
    Metadata    map[string]any // Additional error context
    httpStatus  int            // HTTP status code (internal)
    UserMessage string         // User-friendly message
}
```

## Creating Errors

### Using Predefined Errors
```go
// Basic usage
return apperror.ErrNotFound

// With custom message
return apperror.ErrWrongInput.WithMessage("email cannot be empty")

// With user-friendly message
return apperror.ErrPermission.WithUserMessage("You don't have permission to delete this review")
```

### Creating Custom Errors
```go
// Using helper functions
err := apperror.NewValidationError(4005, "invalid email format")
err := apperror.NewDomainError(1005, "course not found")
err := apperror.NewInfrastructureError(5004, "service unavailable")

// Wrapping existing errors
return apperror.WrapDB(dbErr)
return apperror.WrapInternal(internalErr)
```

### Adding Context
```go
// Add metadata
err := apperror.ErrDB.
    Wrap(dbError).
    WithMetadata("operation", "user_query").
    WithMetadata("user_id", userID).
    WithMetadata("timestamp", time.Now())

// Add context from request
err := apperror.ErrPermission.WithContext(ctx)
```

## Error Handling Patterns

### 1. Repository Layer
```go
func (r *userRepository) GetByID(ctx context.Context, id int) (*User, error) {
    user, err := r.db.User.First(id).Result()
    if err != nil {
        return nil, apperror.WrapDB(err).WithMetadata("operation", "get_user_by_id")
    }
    if user == nil {
        return nil, apperror.ErrNotFound.WithMetadata("user_id", id)
    }
    return user, nil
}
```

### 2. Service Layer
```go
func (s *userService) UpdateUser(ctx context.Context, cmd UpdateUserCommand) error {
    // Validate input
    if cmd.Email == "" {
        return apperror.ErrRequiredField.WithMessage("email is required")
    }
    
    // Check if user exists
    user, err := s.repo.GetByID(ctx, cmd.UserID)
    if err != nil {
        return err // Already wrapped with context
    }
    if user == nil {
        return apperror.ErrNotFound.WithMetadata("user_id", cmd.UserID)
    }
    
    // Update user
    if err := s.repo.Update(ctx, user); err != nil {
        return apperror.WrapDB(err).WithMetadata("operation", "update_user")
    }
    
    return nil
}
```

### 3. HTTP Handlers
```go
func (c *UserController) UpdateUser(ctx *gin.Context) {
    var cmd UpdateUserCommand
    if err := ctx.ShouldBindJSON(&cmd); err != nil {
        web.HandleValidationError(ctx, "invalid request body")
        return
    }
    
    commonCtx := web.GetCommonContext(ctx)
    
    err := c.userService.UpdateUser(commonCtx, &cmd)
    if err != nil {
        web.HandleError(ctx, err)
        return
    }
    
    web.HandleSuccess(ctx, nil)
}
```

## Error Response Format

The system provides consistent error responses:

```json
{
    "success": false,
    "code": 3001,
    "message": "permission denied",
    "category": 3000,
    "user_message": "You don't have permission to perform this action",
    "metadata": {
        "operation": "delete_review",
        "user_id": 123,
        "review_id": 456
    },
    "timestamp": "2024-01-15T10:30:00Z"
}
```

## Error Checking and Matching

### Check for specific error types
```go
if errors.Is(err, apperror.ErrNotFound) {
    // Handle not found error
}

if errors.Is(err, apperror.ErrPermission) {
    // Handle permission error
}
```

### Check error categories
```go
var appErr *apperror.AppError
if errors.As(err, &appErr) {
    switch appErr.Category {
    case apperror.CategoryValidation:
        // Handle validation error
    case apperror.CategoryAuth:
        // Handle auth error
    case apperror.CategoryPermission:
        // Handle permission error
    }
}
```

### Check if error is retryable
```go
var appErr *apperror.AppError
if errors.As(err, &appErr) && appErr.IsRetryable() {
    // Implement retry logic
}
```

## Error Severity and Logging

### Get error severity
```go
var appErr *apperror.AppError
if errors.As(err, &appErr) {
    severity := appErr.GetSeverity() // "low", "medium", "high", "critical"
}
```

### Automatic error handling
```go
// Use the error handler for consistent processing
response := apperror.HandleError(ctx, err)
fmt.Printf("Error response: %+v\n", response)
```

## Best Practices

1. **Always wrap errors**: Use `WrapDB()`, `WrapInternal()` or `.Wrap()` for infrastructure errors
2. **Add context**: Use `.WithMetadata()` to add operation context
3. **User-friendly messages**: Use `.WithUserMessage()` for user-facing messages
4. **Consistent categorization**: Use appropriate error categories
5. **Don't expose sensitive data**: User messages should not contain internal details
6. **Handle errors at boundaries**: Convert errors at repository/service boundaries

## Migration Guide

### Old Pattern
```go
func (s *service) Operation() error {
    if err := s.repo.DoSomething(); err != nil {
        return err // Raw error
    }
    return nil
}
```

### New Pattern
```go
func (s *service) Operation() error {
    if err := s.repo.DoSomething(); err != nil {
        return apperror.WrapDB(err).WithMetadata("operation", "service_operation")
    }
    return nil
}
```

## Testing Error Handling

```go
func TestErrorHandling(t *testing.T) {
    err := someOperation()
    
    // Test for specific error
    if errors.Is(err, apperror.ErrNotFound) {
        t.Log("Correctly identified not found error")
    }
    
    // Test error category
    var appErr *apperror.AppError
    if errors.As(err, &appErr) && appErr.Category == apperror.CategoryValidation {
        t.Log("Correctly identified validation error")
    }
    
    // Test metadata
    if appErr != nil {
        if operation, ok := appErr.Metadata["operation"]; ok {
            t.Logf("Operation: %v", operation)
        }
    }
}
```

## Common Error Scenarios

### Database Errors
```go
// Always wrap database errors
return apperror.WrapDB(err).WithMetadata("table", "users").WithMetadata("operation", "insert")
```

### Validation Errors
```go
// Use validation errors with user-friendly messages
return apperror.ErrInvalidEmail.WithUserMessage("Please enter a valid email address")
```

### Permission Errors
```go
// Include context about what was being accessed
return apperror.ErrPermission.
    WithMessage("user cannot delete this review").
    WithMetadata("user_id", currentUser.ID).
    WithMetadata("review_id", reviewID).
    WithMetadata("owner_id", review.UserID)
```

### Authentication Errors
```go
// Use consistent auth error messages
return apperror.ErrWrongAuth.WithUserMessage("Invalid email or password")
```

This systematic approach ensures consistent error handling throughout the application while providing rich context for debugging and user-friendly messages for end users.