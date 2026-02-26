package apperror

type IAppError interface {
	// Extend default error
	error
	// Make IAppErrors wrapable
	Unwrap() error
	// IAppErrors are mapped to http status codes
	HttpStatus() int
	// IAppErrors provide an ErrorCode for more detailed error message mapping
	ErrorCode() string
}
