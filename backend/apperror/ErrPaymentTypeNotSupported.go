package apperror

import (
	"fmt"
	"net/http"
)

type ErrPaymentTypeNotSupported struct {
	cause   error
	message string
}

// ErrorCode implements [IAppError].
func (e *ErrPaymentTypeNotSupported) ErrorCode() string {
	return "USER_NOT_FOUND"
}

// HttpStatus implements [IAppError].
func (e *ErrPaymentTypeNotSupported) HttpStatus() int {
	return http.StatusNotFound
}

// Unwrap implements [IAppError].
func (e *ErrPaymentTypeNotSupported) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrPaymentTypeNotSupported) Error() string {
	return fmt.Sprintf("ErrPaymentTypeNotSupported: %s - %v", e.message, e.cause)
}

func NewErrPaymentTypeNotSupported(cause error, format string, args ...any) IAppError {
	return &ErrPaymentTypeNotSupported{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
