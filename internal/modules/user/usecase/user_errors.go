package usecase

import "errors"

var (
	ErrInvalidUserID                  = errors.New("users.application: invalid user id")
	ErrNoFieldsToUpdate               = errors.New("users.application: no fields to update")
	ErrCreateActivationPayloadMarshal = errors.New("users.application: create activation payload marshal")
	ErrCreateActivationEventPublish   = errors.New("users.application: create activation event publish")
	ErrSendActivationEmail            = errors.New("users.application: send activation email")
)
