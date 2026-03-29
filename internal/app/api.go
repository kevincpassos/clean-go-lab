package app

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewAPIServer(container *Container) *http.Server {
	router := chi.NewRouter()

	// Global middlewares for observability and resilience.
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	registerAPIRoutes(router, container)

	return &http.Server{
		Addr:    container.Config.HTTPAddr,
		Handler: router,
	}
}
