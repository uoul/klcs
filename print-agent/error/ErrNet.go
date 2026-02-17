package error

import "fmt"

type ErrNet string

// Error implements error.
func (e ErrNet) Error() string {
	return fmt.Sprintf("ErrNet: %s", string(e))
}

func NewErrNet(format string, args ...any) error {
	return ErrNet(fmt.Sprintf(format, args...))
}
