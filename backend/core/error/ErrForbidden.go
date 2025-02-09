package error

import "fmt"

type ErrForbidden string

// Error implements error.
func (e ErrForbidden) Error() string {
	return fmt.Sprintf("ErrForbidden: %s", string(e))
}

func NewErrForbidden(format string, args ...any) error {
	return ErrForbidden(fmt.Sprintf(format, args...))
}
