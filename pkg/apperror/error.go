package apperror

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
)

var (
	ErrDB = &AppError{Code: 3001, Message: "db error"}
)
