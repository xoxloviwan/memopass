package handlers

import (
	"encoding/json"
	"io"
	"iwakho/gopherkeep/internal/model"
	"iwakho/gopherkeep/internal/srv/jwt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var err error
	var u *model.User
	username, password, ok := r.BasicAuth()
	if ok {
		u, err = h.store.GetUser(r.Context(), username)
		if err == nil {
			err = bcrypt.CompareHashAndPassword(u.Hash, []byte(password))
		}
	}
	if !ok || err != nil {
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		h.logger.Error("Auth error", "error", err, "request_id", w.Header().Get("X-Request-ID"))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	err = addAuthData(*u, w)
	if err != nil {
		h.ErrorWithLog(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	creds := model.Creds{}
	user := model.User{}
	// читаем тело запроса
	data, err := io.ReadAll(r.Body)
	if err != nil {
		h.ErrorWithLog(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(data, &creds)
	if err != nil {
		h.ErrorWithLog(w, err.Error(), http.StatusBadRequest)
		return
	}
	if creds.User == "" || creds.Pwd == "" {
		h.ErrorWithLog(w, "empty login or password", http.StatusBadRequest)
		return
	}
	user.Name = creds.User
	user.Hash, err = bcrypt.GenerateFromPassword([]byte(creds.Pwd), 0)
	if err != nil {
		h.ErrorWithLog(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.ID, err = h.store.NewUser(r.Context(), user)
	if err != nil {
		h.ErrorWithLog(w, err.Error(), http.StatusConflict)
		return
	}
	err = addAuthData(user, w)
	if err != nil {
		h.ErrorWithLog(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func addAuthData(user model.User, w http.ResponseWriter) error {
	tkn, err := jwt.BuildJWT(user.Name, user.ID)
	if err != nil {
		return err
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
	return nil
}
