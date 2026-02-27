package apperror

import (
	"fmt"
	"net/http"
)

type ErrNotEnoughBalance struct {
	cause   error
	message string
}

// HttpStatus implements [IAppError].
func (e *ErrNotEnoughBalance) HttpStatus() int {
	return http.StatusNotFound
}

// Unwrap implements [IAppError].
func (e *ErrNotEnoughBalance) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrNotEnoughBalance) Error() string {
	return fmt.Sprintf("ErrNotEnoughBalance: %s - %v", e.message, e.cause)
}

func NewErrNotEnoughBalance(cause error, format string, args ...any) IAppError {
	return &ErrNotEnoughBalance{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
