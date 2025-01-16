package middleware

import "net/http"

const maxMemory = 5 << 20 // 5 MB

func ParseForm(logger logger) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseMultipartForm(maxMemory); err != nil {
				logger.Error("error parsing form", "err", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
