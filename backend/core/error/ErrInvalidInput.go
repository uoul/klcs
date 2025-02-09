package error

import "fmt"

type ErrInvalidInput string

// Error implements error.
func (e ErrInvalidInput) Error() string {
	return fmt.Sprintf("ErrInvalidInput: %s", string(e))
}

func NewErrInvalidInput(format string, args ...any) error {
	return ErrInvalidInput(fmt.Sprintf(format, args...))
}
