package errors

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrAdminEmailExists   = errors.New("admin with given email already exists")
)
