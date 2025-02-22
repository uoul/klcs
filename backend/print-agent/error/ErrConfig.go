package error

import "fmt"

type ErrConfig string

// Error implements error.
func (e ErrConfig) Error() string {
	return fmt.Sprintf("ErrConfig: %s", string(e))
}

func NewErrConfig(format string, args ...any) error {
	return ErrConfig(fmt.Sprintf(format, args...))
}
