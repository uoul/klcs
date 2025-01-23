package error

type ValidationError struct {
	err error
}

// Error implements error.
func (e *ValidationError) Error() string {
	return e.err.Error()
}

func NewValidationError(err error) error {
	return &ValidationError{
		err: err,
	}
}
