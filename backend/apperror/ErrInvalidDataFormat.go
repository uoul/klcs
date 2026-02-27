package apperror

import (
	"fmt"
	"net/http"
)

type ErrInvalidDataFormat struct {
	cause   error
	message string
}

// HttpStatus implements [IAppError].
func (e *ErrInvalidDataFormat) HttpStatus() int {
	return http.StatusServiceUnavailable
}

// Unwrap implements [IAppError].
func (e *ErrInvalidDataFormat) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrInvalidDataFormat) Error() string {
	return fmt.Sprintf("ErrInvalidDataFormat: %s - %v", e.message, e.cause)
}

func NewErrInvalidDataFormat(cause error, format string, args ...any) IAppError {
	return &ErrInvalidDataFormat{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
