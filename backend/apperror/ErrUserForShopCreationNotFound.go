package apperror

import (
	"fmt"
	"net/http"
)

type ErrUserForShopCreationNotFound struct {
	cause   error
	message string
}

// ErrorCode implements [IAppError].
func (e *ErrUserForShopCreationNotFound) ErrorCode() string {
	return "CANNOT_CLOSE_LOCKED_ACCOUNT"
}

// HttpStatus implements [IAppError].
func (e *ErrUserForShopCreationNotFound) HttpStatus() int {
	return http.StatusBadRequest
}

// Unwrap implements [IAppError].
func (e *ErrUserForShopCreationNotFound) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrUserForShopCreationNotFound) Error() string {
	return fmt.Sprintf("ErrUserForShopCreationNotFound: %s - %v", e.message, e.cause)
}

func NewErrUserForShopCreationNotFound(cause error, format string, args ...any) IAppError {
	return &ErrUserForShopCreationNotFound{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
