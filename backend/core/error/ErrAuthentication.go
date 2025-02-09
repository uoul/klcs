package error

import "fmt"

type ErrAuthentication string

// Error implements error.
func (e ErrAuthentication) Error() string {
	return fmt.Sprintf("ErrAuthentication: %s", string(e))
}

func NewErrAuthentication(format string, args ...any) error {
	return ErrAuthentication(fmt.Sprintf(format, args...))
}
