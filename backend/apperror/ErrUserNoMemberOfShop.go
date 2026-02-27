package apperror

import (
	"fmt"
	"net/http"
)

type ErrUserNoMemberOfShop struct {
	cause   error
	message string
}

// HttpStatus implements [IAppError].
func (e *ErrUserNoMemberOfShop) HttpStatus() int {
	return http.StatusBadRequest
}

// Unwrap implements [IAppError].
func (e *ErrUserNoMemberOfShop) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrUserNoMemberOfShop) Error() string {
	return fmt.Sprintf("ErrUserNoMemberOfShop: %s - %v", e.message, e.cause)
}

func NewErrUserNoMemberOfShop(cause error, format string, args ...any) IAppError {
	return &ErrUserNoMemberOfShop{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
