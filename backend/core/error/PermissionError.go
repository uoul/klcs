package error

type PermissionError struct {
	err error
}

// Error implements error.
func (e *PermissionError) Error() string {
	return e.err.Error()
}

func NewPermissionError(err error) error {
	return &PermissionError{
		err: err,
	}
}
