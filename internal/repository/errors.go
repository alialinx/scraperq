package repository

import "errors"

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrEmailExists   = errors.New("email already exists")
	ErrInvalidAPIKey = errors.New("invalid api key")
)
