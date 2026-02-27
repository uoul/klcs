package apperror

import (
	"fmt"
	"net/http"
)

type ErrNoPrinterShopRelation struct {
	cause   error
	message string
}

// HttpStatus implements [IAppError].
func (e *ErrNoPrinterShopRelation) HttpStatus() int {
	return http.StatusNotFound
}

// Unwrap implements [IAppError].
func (e *ErrNoPrinterShopRelation) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e ErrNoPrinterShopRelation) Error() string {
	return fmt.Sprintf("ErrNoPrinterShopRelation: %s - %v", e.message, e.cause)
}

func NewErrNoPrinterShopRelation(cause error, format string, args ...any) IAppError {
	return &ErrNoPrinterShopRelation{
		cause:   cause,
		message: fmt.Sprintf(format, args...),
	}
}
