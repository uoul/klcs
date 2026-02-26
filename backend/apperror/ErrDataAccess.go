package apperror

import (
	"fmt"
	"net/http"
)

type ErrDataAccess struct {
	cause   error
	message string
}

// ErrorCode implements [IAppError].
func (e *ErrDataAccess) ErrorCode() string {
	return "DATA_ACCESS_FAILED"
}

// HttpStatus implements [IAppError].
func (e *ErrDataAccess) HttpStatus() int {
	return http.StatusServiceUnavailable
}

// Unwrap implements [IAppError].
func (e *ErrDataAccess) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrDataAccess) Error() string {
	return fmt.Sprintf("ErrDataAccess: %s - %v", e.message, e.cause)
}

func NewErrDataAccess(cause error, format string, args ...any) IAppError {
	return &ErrDataAccess{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
