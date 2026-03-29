package app

import (
	"net/http"

	userhttp "golab/internal/modules/user/delivery/http"

	"github.com/go-chi/chi/v5"
)

func registerAPIRoutes(router chi.Router, container *Container) {
	router.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	router.Route("/v1", func(r chi.Router) {
		registerV1Routes(r, container)
	})
}

func registerV1Routes(r chi.Router, container *Container) {
	userhttp.RegisterRoutes(r, container.User.HTTPHandler)
}
