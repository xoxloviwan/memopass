package middleware

import (
	"context"
	"fmt"
	"iwakho/gopherkeep/internal/srv/jwt"
	"net/http"
	"strings"
)

type userCtxKey struct{}

var errorAuthHeader = fmt.Sprintf(`%s realm="restricted", error="invalid_token"`, jwt.Bearer)

func CheckAuthOrDirectTo(loginRoute string) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth != "" && strings.HasPrefix(auth, jwt.Bearer) {
				auth = strings.Replace(auth, jwt.Bearer, "", 1)
				user, err := jwt.GetUser(auth)
				if err == nil {
					ctx := context.WithValue(r.Context(), userCtxKey{}, *user)
					r = r.WithContext(ctx)
					next.ServeHTTP(w, r)
					return
				}
			}
			w.Header().Set("WWW-Authenticate", errorAuthHeader)
			w.WriteHeader(http.StatusUnauthorized)
		})
	}
}
