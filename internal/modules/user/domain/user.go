package domain

import (
	"net/mail"
	"strings"
	"time"
)

type User struct {
	ID        int64
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(name, email string) (*User, error) {
	if err := ValidateName(&name); err != nil {
		return nil, err
	}

	if err := ValidateEmail(&email); err != nil {
		return nil, err
	}

	return &User{
		Name:  name,
		Email: email,
	}, nil
}

func ValidateName(name *string) error {
	*name = strings.TrimSpace(*name)
	if *name == "" {
		return ErrNameRequired
	}
	return nil
}

func ValidateEmail(email *string) error {
	*email = strings.TrimSpace(strings.ToLower(*email))
	if _, err := mail.ParseAddress(*email); err != nil {
		return ErrInvalidEmail
	}
	return nil
}
