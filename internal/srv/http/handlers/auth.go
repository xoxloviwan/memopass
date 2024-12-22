package handlers

import (
	"fmt"
	"iwakho/gopherkeep/internal/srv/jwt"
	"net/http"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		// If the Authentication header is not present, is invalid, or the
		// username or password is wrong, then set a WWW-Authenticate
		// header to inform the client that we expect them to use basic
		// authentication and send a 401 Unauthorized response.
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	fmt.Println(username, password) // TODO найти пользователя в базе данных и проверить его пароль

	// Сгенерируем новый токен доступа
	tkn, err := jwt.BuildJWT(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", jwt.Bearer+tkn)
	cookie := &http.Cookie{
		Name:     "tkn",
		Value:    tkn,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
