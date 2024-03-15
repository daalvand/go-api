package repositories

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrDuplicateEmail    = errors.New("user with this email already exists")
	ErrIncorrectPassword = errors.New("user or password are incorrect")
)
