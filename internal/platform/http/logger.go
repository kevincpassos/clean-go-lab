package platformhttp

import (
	"log/slog"
	"net/http"
)

func RequestLogger(logger *slog.Logger, r *http.Request) *slog.Logger {
	if logger == nil {
		logger = slog.Default()
	}

	return logger.With(
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
	)
}
