package apperror

import (
	"fmt"
	"net/http"
)

type ErrUnauthorized struct {
	cause   error
	message string
}

// ErrorCode implements [IAppError].
func (e *ErrUnauthorized) ErrorCode() string {
	return "NOT_AUTHENTICATED"
}

// HttpStatus implements [IAppError].
func (e *ErrUnauthorized) HttpStatus() int {
	return http.StatusUnauthorized
}

// Unwrap implements [IAppError].
func (e *ErrUnauthorized) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrUnauthorized) Error() string {
	return fmt.Sprintf("ErrUnauthorized: %s - %v", e.message, e.cause)
}

func NewErrUnauthorized(cause error, format string, args ...any) IAppError {
	return &ErrUnauthorized{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
