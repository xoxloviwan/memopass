package middleware

import (
	"log/slog"
	"net/http"

	"github.com/felixge/httpsnoop"
	"github.com/google/uuid"
)

type logger interface {
	Error(msg string, args ...any)
	Info(msg string, args ...any)
	With(args ...any) *slog.Logger
}

type middleware func(http.Handler) http.Handler

func Logging(logger logger) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := uuid.New().String()
			w.Header().Set("X-Request-ID", reqID)
			logger := logger.With(
				slog.String("request_id", reqID),
			)
			logger.Info(
				"REQ",
				slog.String("method", r.Method),
				slog.String("uri", r.URL.String()),
				slog.String("ip", r.RemoteAddr),
				slog.String("size", r.Header.Get("Content-Length")),
				slog.String("user_agent", r.Header.Get("User-Agent")),
			)

			// this runs handler next and captures information about HTTP request
			m := httpsnoop.CaptureMetrics(next, w, r)
			logger.Info(
				"RES",
				slog.Int("status", m.Code),
				slog.Duration("duration", m.Duration),
				slog.Int64("size", m.Written),
			)
		})
	}
}
