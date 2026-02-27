package apperror

type IAppError interface {
	// Extend default error
	error
	// Make IAppErrors wrapable
	Unwrap() error
	// IAppErrors are mapped to http status codes
	HttpStatus() int
}
