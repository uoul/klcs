package error

type NotFoundError struct {
	err error
}

// Error implements error.
func (e *NotFoundError) Error() string {
	return e.err.Error()
}

func NewNotFoundError(err error) error {
	return &NotFoundError{
		err: err,
	}
}
