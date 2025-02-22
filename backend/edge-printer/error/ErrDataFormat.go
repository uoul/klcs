package error

import "fmt"

type ErrDataFormat string

// Error implements error.
func (e ErrDataFormat) Error() string {
	return fmt.Sprintf("ErrDataFormat: %s", string(e))
}

func NewErrDataFormat(format string, args ...any) error {
	return ErrDataFormat(fmt.Sprintf(format, args...))
}
