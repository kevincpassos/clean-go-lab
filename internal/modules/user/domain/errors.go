package domain

import "errors"

var (
	ErrNameRequired = errors.New("users.domain: name is required")
	ErrInvalidEmail = errors.New("users.domain: invalid email")
)
