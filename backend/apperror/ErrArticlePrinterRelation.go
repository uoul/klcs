package apperror

import (
	"fmt"
	"net/http"
)

type ErrArticlePrinterRelation struct {
	cause   error
	message string
}

// ErrorCode implements [IAppError].
func (e *ErrArticlePrinterRelation) ErrorCode() string {
	return "NO_PRINTER_SHOP_RELATION"
}

// HttpStatus implements [IAppError].
func (e *ErrArticlePrinterRelation) HttpStatus() int {
	return http.StatusNotFound
}

// Unwrap implements [IAppError].
func (e *ErrArticlePrinterRelation) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrArticlePrinterRelation) Error() string {
	return fmt.Sprintf("ErrArticlePrinterRelation: %s - %v", e.message, e.cause)
}

func NewErrArticlePrinterRelation(cause error, format string, args ...any) IAppError {
	return &ErrArticlePrinterRelation{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
