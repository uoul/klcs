package apperror

import (
	"fmt"
	"net/http"
)

type ErrUserShopPermission struct {
	cause   error
	message string
}

// HttpStatus implements [IAppError].
func (e *ErrUserShopPermission) HttpStatus() int {
	return http.StatusUnauthorized
}

// Unwrap implements [IAppError].
func (e *ErrUserShopPermission) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrUserShopPermission) Error() string {
	return fmt.Sprintf("ErrUserShopPermission: %s - %v", e.message, e.cause)
}

func NewErrUserShopPermission(cause error, format string, args ...any) IAppError {
	return &ErrUserShopPermission{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
