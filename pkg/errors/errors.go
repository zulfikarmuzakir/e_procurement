package errors

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInternalServer     = errors.New("internal server error")
	ErrUserNotActive      = errors.New("user not active")
)

type AppError struct {
	Err     error
	Message string
	Code    int
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func NewAppError(err error, message string, code int) *AppError {
	return &AppError{Err: err, Message: message, Code: code}
}
