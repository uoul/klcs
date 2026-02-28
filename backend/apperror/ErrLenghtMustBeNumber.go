package apperror

import (
	"fmt"
	"net/http"
)

type ErrLengthMustBeNumber struct {
	cause   error
	message string
}

// HttpStatus implements [IAppError].
func (e *ErrLengthMustBeNumber) HttpStatus() int {
	return http.StatusServiceUnavailable
}

// Unwrap implements [IAppError].
func (e *ErrLengthMustBeNumber) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrLengthMustBeNumber) Error() string {
	return fmt.Sprintf("ErrLengthMustBeNumber: %s - %v", e.message, e.cause)
}

func NewErrLengthMustBeNumber(cause error, format string, args ...any) IAppError {
	return &ErrLengthMustBeNumber{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
