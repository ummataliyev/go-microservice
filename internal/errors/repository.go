package errors

import "errors"

var (
	ErrNotFound     = errors.New("resource not found")
	ErrCannotCreate = errors.New("cannot create resource")
	ErrCannotUpdate = errors.New("cannot update resource")
	ErrCannotDelete = errors.New("cannot delete resource")
	ErrInvalidInput = errors.New("invalid input")
)
