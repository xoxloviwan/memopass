package middleware

import (
	"context"
	"fmt"
	"iwakho/gopherkeep/internal/model"
	"iwakho/gopherkeep/internal/srv/jwt"
	"net/http"
	"strings"
)

var errorAuthHeader = fmt.Sprintf(`%s realm="restricted", error="invalid_token"`, jwt.Bearer)

func (m *Middlewares) CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			cookie, err := r.Cookie("tkn")
			if err == nil {
				auth = jwt.Bearer + cookie.Value
			}
		}
		if auth != "" && strings.HasPrefix(auth, jwt.Bearer) {
			auth = strings.Replace(auth, jwt.Bearer, "", 1)
			user, err := jwt.GetUser(auth)
			if err == nil {
				ctx := context.WithValue(r.Context(), model.UserIDCtxKey{}, *user)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}
		}
		w.Header().Set("WWW-Authenticate", errorAuthHeader)
		w.WriteHeader(http.StatusUnauthorized)
	})
}
