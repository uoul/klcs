package apperror

import (
	"fmt"
	"net/http"
)

type ErrAccountNotFound struct {
	cause   error
	message string
}

// ErrorCode implements [IAppError].
func (e *ErrAccountNotFound) ErrorCode() string {
	return "ACCOUNT_NOT_FOUND"
}

// HttpStatus implements [IAppError].
func (e *ErrAccountNotFound) HttpStatus() int {
	return http.StatusNotFound
}

// Unwrap implements [IAppError].
func (e *ErrAccountNotFound) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrAccountNotFound) Error() string {
	return fmt.Sprintf("ErrAccountNotFound: %s - %v", e.message, e.cause)
}

func NewErrAccountNotFound(cause error, format string, args ...any) IAppError {
	return &ErrAccountNotFound{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
