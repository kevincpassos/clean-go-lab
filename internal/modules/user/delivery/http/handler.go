package http

import (
	"golab/internal/modules/user/usecase"
	platformhttp "golab/internal/platform/http"
	"log/slog"
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	validate *validator.Validate
	logger   *slog.Logger
	useCase  *usecase.UserUseCase
}

func NewHandler(useCase *usecase.UserUseCase, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}

	return &Handler{
		validate: validator.New(),
		logger:   logger,
		useCase:  useCase,
	}
}

func (h *Handler) requestLogger(r *nethttp.Request) *slog.Logger {
	return platformhttp.RequestLogger(h.logger, r)
}

func (h *Handler) Create(w nethttp.ResponseWriter, r *nethttp.Request) {
	log := h.requestLogger(r)

	var req CreateUserRequest
	if err := platformhttp.DecodeJSONBody(r, &req); err != nil {
		log.Warn("http create user invalid body", slog.Any("error", err))
		writeError(w, platformhttp.ErrInvalidRequestBody)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		log.Warn("http create user validation failed", slog.Any("error", err))
		writeError(w, err)
		return
	}

	out, err := h.useCase.Create(r.Context(), usecase.CreateUserInput{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		log.Error("http create user failed", slog.Any("error", err))
		writeError(w, err)
		return
	}

	platformhttp.WriteJSON(w, nethttp.StatusCreated, toResponse(out))
}

func (h *Handler) Patch(w nethttp.ResponseWriter, r *nethttp.Request) {
	log := h.requestLogger(r)

	id, err := extractIDParam(r)
	if err != nil {
		log.Warn("http patch user invalid id", slog.Any("error", err))
		writeError(w, usecase.ErrInvalidUserID)
		return
	}

	log = log.With(slog.Int64("user_id", id))

	var req PatchUserRequest
	if err := platformhttp.DecodeJSONBody(r, &req); err != nil {
		log.Warn("http patch user invalid body", slog.Any("error", err))
		writeError(w, platformhttp.ErrInvalidRequestBody)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		log.Warn("http patch user validation failed", slog.Any("error", err))
		writeError(w, err)
		return
	}

	out, err := h.useCase.Patch(r.Context(), usecase.PatchUserInput{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		log.Error("http patch user failed", slog.Any("error", err))
		writeError(w, err)
		return
	}

	platformhttp.WriteJSON(w, nethttp.StatusOK, toResponse(out))
}

func (h *Handler) Delete(w nethttp.ResponseWriter, r *nethttp.Request) {
	log := h.requestLogger(r)

	id, err := extractIDParam(r)
	if err != nil {
		log.Warn("http delete user invalid id", slog.Any("error", err))
		writeError(w, usecase.ErrInvalidUserID)
		return
	}

	log = log.With(slog.Int64("user_id", id))

	if err := h.useCase.Delete(r.Context(), usecase.DeleteUserInput{ID: id}); err != nil {
		log.Error("http delete user failed", slog.Any("error", err))
		writeError(w, err)
		return
	}

	w.WriteHeader(nethttp.StatusNoContent)
}

func writeError(w nethttp.ResponseWriter, err error) {
	platformhttp.WriteError(w, err, UserMapError)
}

func extractIDParam(r *nethttp.Request) (int64, error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return 0, usecase.ErrInvalidUserID
	}

	parsedID, err := parseUserID(id)
	if err != nil {
		return 0, err
	}

	return parsedID, nil
}
