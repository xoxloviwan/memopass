package handlers

import "net/http"

func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
