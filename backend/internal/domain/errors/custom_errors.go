package errors

import "errors"

// CustomError representa um erro estruturado com código, mensagem e classificação por tipo.
type CustomError struct {
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	errorType ErrorType
	details   []any
}

func NewCustomError(code, message string, errorType ErrorType, details []any) CustomError {
	return CustomError{
		Code:      code,
		Message:   message,
		errorType: errorType,
		details:   details,
	}
}

func (e CustomError) Error() string {
	return e.Message
}

func (e CustomError) Type() ErrorType {
	return e.errorType
}

func (e CustomError) Details() []any {
	return e.details
}

func (e CustomError) Is(target error) bool {
	t, ok := target.(*CustomError)
	return ok && e.errorType == t.errorType
}

func (e *CustomError) As(target any) bool {
	if t, ok := target.(*CustomError); ok {
		*t = *e
		return true
	}
	return false
}

func NewValidationError(code, message string, details []any) CustomError {
	return NewCustomError(code, message, ErrorTypeValidation, details)
}

func NewBadRequestError(code, message string) CustomError {
	return NewCustomError(code, message, ErrorTypeBadRequest, nil)
}

func NewNotFoundError(message string) CustomError {
	return NewCustomError("not_found", message, ErrorTypeNotFound, nil)
}

func NewAlreadyExistsError(message string) CustomError {
	return NewCustomError("already_exists", message, ErrorTypeAlreadyExists, nil)
}

func IsNotFoundError(err error) bool {
	return isOfType(err, ErrorTypeNotFound)
}

func IsValidationError(err error) bool {
	return isOfType(err, ErrorTypeValidation)
}

func IsAlreadyExistsError(err error) bool {
	return isOfType(err, ErrorTypeAlreadyExists)
}

// IsCustomError tenta extrair um CustomError de err via errors.As.
// Retorna true e preenche target se bem-sucedido.
func IsCustomError(err error, target *CustomError) bool {
	return errors.As(err, target)
}

func isOfType(err error, t ErrorType) bool {
	ce := &CustomError{}
	return errors.As(err, &ce) && ce.Type() == t
}
