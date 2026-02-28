package apperror

import (
	"fmt"
	"net/http"
)

type ErrArticleShopRelation struct {
	cause   error
	message string
}

// HttpStatus implements [IAppError].
func (e *ErrArticleShopRelation) HttpStatus() int {
	return http.StatusNotFound
}

// Unwrap implements [IAppError].
func (e *ErrArticleShopRelation) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrArticleShopRelation) Error() string {
	return fmt.Sprintf("ErrArticleShopRelation: %s - %v", e.message, e.cause)
}

func NewErrArticleShopRelation(cause error, format string, args ...any) IAppError {
	return &ErrArticleShopRelation{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
