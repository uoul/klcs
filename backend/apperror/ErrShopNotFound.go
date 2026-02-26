package apperror

import (
	"fmt"
	"net/http"
)

type ErrShopNotFound struct {
	cause   error
	message string
}

// ErrorCode implements [IAppError].
func (e *ErrShopNotFound) ErrorCode() string {
	return "USER_NOT_FOUND"
}

// HttpStatus implements [IAppError].
func (e *ErrShopNotFound) HttpStatus() int {
	return http.StatusNotFound
}

// Unwrap implements [IAppError].
func (e *ErrShopNotFound) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrShopNotFound) Error() string {
	return fmt.Sprintf("ErrShopNotFound: %s - %v", e.message, e.cause)
}

func NewErrShopNotFound(cause error, format string, args ...any) IAppError {
	return &ErrShopNotFound{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
