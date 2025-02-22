package error

import "fmt"

type ErrIO string

// Error implements error.
func (e ErrIO) Error() string {
	return fmt.Sprintf("ErrIO: %s", string(e))
}

func NewErrIO(format string, args ...any) error {
	return ErrIO(fmt.Sprintf(format, args...))
}
