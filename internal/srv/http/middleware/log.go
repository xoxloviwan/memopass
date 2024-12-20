package middleware

import (
	"log/slog"
	"net/http"

	"github.com/felixge/httpsnoop"
	"github.com/google/uuid"
)

type logger interface {
	Info(msg string, args ...any)
}

type middleware func(http.Handler) http.Handler

func Logging(logger logger) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := uuid.New().String()
			logger.Info(
				"REQ",
				slog.String("method", r.Method),
				slog.String("uri", r.URL.String()),
				slog.String("request_id", reqID),
				slog.String("ip", r.RemoteAddr),
				slog.String("user_agent", r.Header.Get("User-Agent")),
			)

			// this runs handler next and captures information about HTTP request
			m := httpsnoop.CaptureMetrics(next, w, r)
			logger.Info(
				"RES",
				slog.Int("status", m.Code),
				slog.Duration("duration", m.Duration),
				slog.Int64("size", m.Written),
				slog.String("request_id", reqID),
			)
		})
	}
}
