package apperror

import (
	"fmt"
	"net/http"
)

type ErrPostLockedAccount struct {
	cause   error
	message string
}

// HttpStatus implements [IAppError].
func (e *ErrPostLockedAccount) HttpStatus() int {
	return http.StatusBadRequest
}

// Unwrap implements [IAppError].
func (e *ErrPostLockedAccount) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrPostLockedAccount) Error() string {
	return fmt.Sprintf("ErrPostLockedAccount: %s - %v", e.message, e.cause)
}

func NewErrPostLockedAccount(cause error, format string, args ...any) IAppError {
	return &ErrPostLockedAccount{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
