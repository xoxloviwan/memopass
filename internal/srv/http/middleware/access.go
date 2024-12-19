package middleware

import "net/http"

func CheckAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// auth := r.Header.Get("Authorization")
		h.ServeHTTP(w, r)
	})
}
