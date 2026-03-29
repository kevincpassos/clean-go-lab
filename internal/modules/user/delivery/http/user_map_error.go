package http

import (
	"errors"
	"golab/internal/modules/user/domain"
	"golab/internal/modules/user/usecase"
	"golab/internal/modules/user/usecase/ports"
	platformhttp "golab/internal/platform/http"
	nethttp "net/http"

	"github.com/go-playground/validator/v10"
)

func UserMapError(err error) (int, string) {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		return nethttp.StatusBadRequest, err.Error()
	}

	switch {
	case errors.Is(err, platformhttp.ErrInvalidRequestBody):
		return nethttp.StatusBadRequest, "invalid request body"
	case errors.Is(err, platformhttp.ErrValidationFailed):
		return nethttp.StatusBadRequest, "validation failed"

	case errors.Is(err, usecase.ErrInvalidUserID):
		return nethttp.StatusBadRequest, "invalid user id"
	case errors.Is(err, usecase.ErrNoFieldsToUpdate):
		return nethttp.StatusBadRequest, "no fields to update"
	case errors.Is(err, usecase.ErrCreateActivationPayloadMarshal):
		return nethttp.StatusInternalServerError, "failed to build activation payload"
	case errors.Is(err, usecase.ErrCreateActivationEventPublish):
		return nethttp.StatusInternalServerError, "failed to publish activation event"
	case errors.Is(err, usecase.ErrSendActivationEmail):
		return nethttp.StatusInternalServerError, "failed to send activation email"

	case errors.Is(err, domain.ErrNameRequired):
		return nethttp.StatusBadRequest, "name is required"
	case errors.Is(err, domain.ErrInvalidEmail):
		return nethttp.StatusBadRequest, "invalid email"

	case errors.Is(err, ports.ErrUserNotFound):
		return nethttp.StatusNotFound, "user not found"
	case errors.Is(err, ports.ErrEmailConflict):
		return nethttp.StatusConflict, "email already exists"

	default:
		return nethttp.StatusInternalServerError, "internal server error"
	}
}
