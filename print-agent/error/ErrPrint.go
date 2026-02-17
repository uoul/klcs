package error

import "fmt"

type ErrPrint string

// Error implements error.
func (e ErrPrint) Error() string {
	return fmt.Sprintf("ErrPrint: %s", string(e))
}

func NewErrPrint(format string, args ...any) error {
	return ErrPrint(fmt.Sprintf(format, args...))
}
