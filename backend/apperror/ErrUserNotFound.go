package apperror

import (
	"fmt"
	"net/http"
)

type ErrUserNotFound struct {
	cause   error
	message string
}

// ErrorCode implements [IAppError].
func (e *ErrUserNotFound) ErrorCode() string {
	return "USER_NOT_FOUND"
}

// HttpStatus implements [IAppError].
func (e *ErrUserNotFound) HttpStatus() int {
	return http.StatusNotFound
}

// Unwrap implements [IAppError].
func (e *ErrUserNotFound) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrUserNotFound) Error() string {
	return fmt.Sprintf("ErrUserNotFound: %s - %v", e.message, e.cause)
}

func NewErrUserNotFound(cause error, format string, args ...any) IAppError {
	return &ErrUserNotFound{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
