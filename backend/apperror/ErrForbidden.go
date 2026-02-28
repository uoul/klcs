package apperror

import (
	"fmt"
	"net/http"
)

type ErrMissingOidcRole struct {
	cause   error
	message string
}

// HttpStatus implements [IAppError].
func (e *ErrMissingOidcRole) HttpStatus() int {
	return http.StatusNotFound
}

// Unwrap implements [IAppError].
func (e *ErrMissingOidcRole) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrMissingOidcRole) Error() string {
	return fmt.Sprintf("ErrMissingOidcRole: %s - %v", e.message, e.cause)
}

func NewErrMissingOidcRole(cause error, format string, args ...any) IAppError {
	return &ErrMissingOidcRole{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
