package apperror

import (
	"fmt"
	"net/http"
)

type ErrArtilceNotFound struct {
	cause   error
	message string
}

// ErrorCode implements [IAppError].
func (e *ErrArtilceNotFound) ErrorCode() string {
	return "USER_NOT_FOUND"
}

// HttpStatus implements [IAppError].
func (e *ErrArtilceNotFound) HttpStatus() int {
	return http.StatusNotFound
}

// Unwrap implements [IAppError].
func (e *ErrArtilceNotFound) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrArtilceNotFound) Error() string {
	return fmt.Sprintf("ErrArtilceNotFound: %s - %v", e.message, e.cause)
}

func NewErrArtilceNotFound(cause error, format string, args ...any) IAppError {
	return &ErrArtilceNotFound{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
