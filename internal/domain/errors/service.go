package errors

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidTokenType   = errors.New("invalid token type")
	ErrLoginLocked        = errors.New("login locked due to too many failed attempts")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrServiceInvalidInput = errors.New("invalid service input")
)
