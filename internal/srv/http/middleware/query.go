package middleware

import (
	"context"
	"iwakho/gopherkeep/internal/model"
	"net/http"
	"strconv"
)

const (
	limitDefault  = 10
	offsetDefault = 0
)

func ParseQueryParams(log logger) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			queries := r.URL.Query()
			limitStr := queries.Get("limit")
			offsetStr := queries.Get("offset")

			limit := limitDefault
			offset := offsetDefault
			var err error
			if limitStr != "" {
				limit, err = strconv.Atoi(limitStr)
				if err != nil {
					log.Error("wrong limit value", "err", err.Error())
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}
			if offsetStr != "" {
				offset, err = strconv.Atoi(offsetStr)
				if err != nil {
					log.Error("wrong offset value", "err", err.Error())
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}
			ctx := context.WithValue(r.Context(), model.LimitCtxKey{}, limit)
			ctx = context.WithValue(ctx, model.OffsetCtxKey{}, offset)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
