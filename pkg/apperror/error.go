package apperror

import "net/http"

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Wrap(err error) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: e.Message,
		Err:     err,
	}
}

func (e *AppError) WithMessage(message string) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: message,
		Err:     e.Err,
	}
}

func (e *AppError) HTTPStatus() int {
	switch e.Code {
	case 1001: // ErrNotFound
		return http.StatusNotFound
	case 1002: // ErrWrongInput
		return http.StatusBadRequest
	case 1003: // ErrRateLimit
		return http.StatusTooManyRequests
	case 2002: // ErrWrongAuth
		return http.StatusUnauthorized
	case 2003: // ErrSession
		return http.StatusUnauthorized
	case 2004: // ErrSuspended
		return http.StatusForbidden
	case 2008: // ErrPermission
		return http.StatusForbidden
	case 3001: // ErrDB
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

var (
	ErrNotFound   = &AppError{Code: 1001, Message: "not found"}
	ErrWrongInput = &AppError{Code: 1002, Message: "wrong input"}
	ErrRateLimit  = &AppError{Code: 1003, Message: "rate limit"}
)

var (
	ErrExpired        = &AppError{Code: 2001, Message: "expired"}
	ErrWrongAuth      = &AppError{Code: 2002, Message: "wrong user email or password"}
	ErrSession        = &AppError{Code: 2003, Message: "session error"}
	ErrSuspended      = &AppError{Code: 2004, Message: "user has been suspended"}
	ErrNoTargetCourse = &AppError{Code: 2005, Message: "no target course"}
	ErrNoSemester     = &AppError{Code: 2006, Message: "no semester"}
	ErrUserNotFound   = &AppError{Code: 2007, Message: "user not found"}
	ErrPermission     = &AppError{Code: 2008, Message: "permission denied"}
)

var (
	ErrDB = &AppError{Code: 3001, Message: "db error"}
)
