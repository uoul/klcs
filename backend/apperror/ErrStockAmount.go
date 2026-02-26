package apperror

import (
	"fmt"
	"net/http"
)

type ErrStockAmount struct {
	cause   error
	message string
}

// ErrorCode implements [IAppError].
func (e *ErrStockAmount) ErrorCode() string {
	return "USER_NOT_FOUND"
}

// HttpStatus implements [IAppError].
func (e *ErrStockAmount) HttpStatus() int {
	return http.StatusNotFound
}

// Unwrap implements [IAppError].
func (e *ErrStockAmount) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrStockAmount) Error() string {
	return fmt.Sprintf("ErrStockAmount: %s - %v", e.message, e.cause)
}

func NewErrStockAmount(cause error, format string, args ...any) IAppError {
	return &ErrStockAmount{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
