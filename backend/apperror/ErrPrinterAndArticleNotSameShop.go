package apperror

import (
	"fmt"
	"net/http"
)

type ErrPrinterAndArticleNotSameShop struct {
	cause   error
	message string
}

// HttpStatus implements [IAppError].
func (e *ErrPrinterAndArticleNotSameShop) HttpStatus() int {
	return http.StatusBadRequest
}

// Unwrap implements [IAppError].
func (e *ErrPrinterAndArticleNotSameShop) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrPrinterAndArticleNotSameShop) Error() string {
	return fmt.Sprintf("ErrPrinterAndArticleNotSameShop: %s - %v", e.message, e.cause)
}

func NewErrPrinterAndArticleNotSameShop(cause error, format string, args ...any) IAppError {
	return &ErrPrinterAndArticleNotSameShop{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
