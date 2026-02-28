package apperror

import (
	"fmt"
	"net/http"
)

type ErrInternal struct {
	cause   error
	message string
}

// HttpStatus implements [IAppError].
func (e *ErrInternal) HttpStatus() int {
	return http.StatusInternalServerError
}

// Unwrap implements [IAppError].
func (e *ErrInternal) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrInternal) Error() string {
	return fmt.Sprintf("ErrInternal: %s - %v", e.message, e.cause)
}

func NewErrInternal(cause error, format string, args ...any) IAppError {
	return &ErrInternal{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
