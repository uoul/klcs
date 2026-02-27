package apperror

import (
	"fmt"
	"net/http"
)

type ErrNoMatchingArticleIds struct {
	cause   error
	message string
}

// HttpStatus implements [IAppError].
func (e *ErrNoMatchingArticleIds) HttpStatus() int {
	return http.StatusNotFound
}

// Unwrap implements [IAppError].
func (e *ErrNoMatchingArticleIds) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrNoMatchingArticleIds) Error() string {
	return fmt.Sprintf("ErrNoMatchingArticleIds: %s - %v", e.message, e.cause)
}

func NewErrNoMatchingArticleIds(cause error, format string, args ...any) IAppError {
	return &ErrNoMatchingArticleIds{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
