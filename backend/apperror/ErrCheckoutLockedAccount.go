package apperror

import (
	"fmt"
	"net/http"
)

type ErrCheckoutLockedAccount struct {
	cause   error
	message string
}

// ErrorCode implements [IAppError].
func (e *ErrCheckoutLockedAccount) ErrorCode() string {
	return "USER_NOT_FOUND"
}

// HttpStatus implements [IAppError].
func (e *ErrCheckoutLockedAccount) HttpStatus() int {
	return http.StatusNotFound
}

// Unwrap implements [IAppError].
func (e *ErrCheckoutLockedAccount) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrCheckoutLockedAccount) Error() string {
	return fmt.Sprintf("ErrCheckoutLockedAccount: %s - %v", e.message, e.cause)
}

func NewErrCheckoutLockedAccount(cause error, format string, args ...any) IAppError {
	return &ErrCheckoutLockedAccount{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
