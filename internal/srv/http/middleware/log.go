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

type Middlewares struct {
	logger
}

func New(logger logger) *Middlewares {
	return &Middlewares{logger}
}

func (mw *Middlewares) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New().String()
		w.Header().Set("X-Request-ID", reqID)
		mw.logger = mw.logger.With(
			slog.String("request_id", reqID),
		)
		logger := mw.logger
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
