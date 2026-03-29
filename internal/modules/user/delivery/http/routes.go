package http

import (
	"golab/internal/modules/user/usecase"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.Post("/users", handler.Create)
	r.Patch("/users/{id}", handler.Patch)
	r.Delete("/users/{id}", handler.Delete)
}

func parseUserID(value string) (int64, error) {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, usecase.ErrInvalidUserID
	}
	return id, nil
}
