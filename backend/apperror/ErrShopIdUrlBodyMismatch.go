package apperror

import (
	"fmt"
	"net/http"
)

type ErrShopIdUrlBodyMismatch struct {
	cause   error
	message string
}

// HttpStatus implements [IAppError].
func (e *ErrShopIdUrlBodyMismatch) HttpStatus() int {
	return http.StatusBadRequest
}

// Unwrap implements [IAppError].
func (e *ErrShopIdUrlBodyMismatch) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrShopIdUrlBodyMismatch) Error() string {
	return fmt.Sprintf("ErrShopIdUrlBodyMismatch: %s - %v", e.message, e.cause)
}

func NewErrShopIdUrlBodyMismatch(cause error, format string, args ...any) IAppError {
	return &ErrShopIdUrlBodyMismatch{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
