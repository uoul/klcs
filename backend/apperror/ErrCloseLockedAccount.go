package apperror

import (
	"fmt"
	"net/http"
)

type ErrCloseLockedAccount struct {
	cause   error
	message string
}

// HttpStatus implements [IAppError].
func (e *ErrCloseLockedAccount) HttpStatus() int {
	return http.StatusBadRequest
}

// Unwrap implements [IAppError].
func (e *ErrCloseLockedAccount) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrCloseLockedAccount) Error() string {
	return fmt.Sprintf("ErrCloseLockedAccount: %s - %v", e.message, e.cause)
}

func NewErrCloseLockedAccount(cause error, format string, args ...any) IAppError {
	return &ErrCloseLockedAccount{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
