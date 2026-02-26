package apperror

import (
	"fmt"
	"net/http"
)

type ErrNoAccountIdForCartPayment struct {
	cause   error
	message string
}

// ErrorCode implements [IAppError].
func (e *ErrNoAccountIdForCartPayment) ErrorCode() string {
	return "USER_NOT_FOUND"
}

// HttpStatus implements [IAppError].
func (e *ErrNoAccountIdForCartPayment) HttpStatus() int {
	return http.StatusNotFound
}

// Unwrap implements [IAppError].
func (e *ErrNoAccountIdForCartPayment) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrNoAccountIdForCartPayment) Error() string {
	return fmt.Sprintf("ErrNoAccountIdForCartPayment: %s - %v", e.message, e.cause)
}

func NewErrNoAccountIdForCartPayment(cause error, format string, args ...any) IAppError {
	return &ErrNoAccountIdForCartPayment{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
