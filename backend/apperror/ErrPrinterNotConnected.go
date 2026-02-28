package apperror

import (
	"fmt"
	"net/http"
)

type ErrPrinterNotConnected struct {
	cause   error
	message string
}

// HttpStatus implements [IAppError].
func (e *ErrPrinterNotConnected) HttpStatus() int {
	return http.StatusNotFound
}

// Unwrap implements [IAppError].
func (e *ErrPrinterNotConnected) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrPrinterNotConnected) Error() string {
	return fmt.Sprintf("ErrPrinterNotConnected: %s - %v", e.message, e.cause)
}

func NewErrPrinterNotConnected(cause error, format string, args ...any) IAppError {
	return &ErrPrinterNotConnected{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
