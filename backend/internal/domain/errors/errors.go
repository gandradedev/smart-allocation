package errors

var (
	ErrorTypeNotFound      = ErrorType{"not-found"}
	ErrorTypeValidation    = ErrorType{"invalid-input"}
	ErrorTypeBadRequest    = ErrorType{"bad-request"}
	ErrorTypeAlreadyExists = ErrorType{"already-exists"}
)

type ErrorType struct {
	t string
}
