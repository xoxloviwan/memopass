package middleware

import (
	"context"
	"iwakho/gopherkeep/internal/model"
	"net/http"
	"strconv"
)

func (m *Middlewares) ParseQueryParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queries := r.URL.Query()
		limitStr := queries.Get("limit")
		offsetStr := queries.Get("offset")

		var limit, offset int
		var err error
		if limitStr != "" {
			limit, err = strconv.Atoi(limitStr)
			if err != nil {
				m.logger.Error("wrong limit value", "err", err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		if offsetStr != "" {
			offset, err = strconv.Atoi(offsetStr)
			if err != nil {
				m.logger.Error("wrong offset value", "err", err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		ctx := context.WithValue(r.Context(), model.LimitCtxKey{}, limit)
		ctx = context.WithValue(ctx, model.OffsetCtxKey{}, offset)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
