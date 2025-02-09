package error

import "fmt"

type ErrValidation string

// Error implements error.
func (e ErrValidation) Error() string {
	return fmt.Sprintf("ErrValidation: %v", string(e))
}

func NewErrValidation(format string, args ...any) error {
	return ErrValidation(fmt.Sprintf(format, args...))
}
