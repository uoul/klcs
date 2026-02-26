package apperror

import (
	"fmt"
	"net/http"
)

type ErrPrinterNotFound struct {
	cause   error
	message string
}

// ErrorCode implements [IAppError].
func (e *ErrPrinterNotFound) ErrorCode() string {
	return "USER_NOT_FOUND"
}

// HttpStatus implements [IAppError].
func (e *ErrPrinterNotFound) HttpStatus() int {
	return http.StatusNotFound
}

// Unwrap implements [IAppError].
func (e *ErrPrinterNotFound) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrPrinterNotFound) Error() string {
	return fmt.Sprintf("ErrPrinterNotFound: %s - %v", e.message, e.cause)
}

func NewErrPrinterNotFound(cause error, format string, args ...any) IAppError {
	return &ErrPrinterNotFound{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
