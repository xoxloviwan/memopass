package handlers

import (
	"encoding/json"
	"iwakho/gopherkeep/internal/model"
	"net/http"
	"time"
)

func (h *Handler) AddPair(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.UserIDCtxKey{}).(int)
	if r.PostForm.Get("login") == "" || r.PostForm.Get("password") == "" {
		h.ErrorWithLog(w, "login or password is empty", http.StatusBadRequest)
		return
	}
	pairs := model.PairInfo{
		Pair: model.Pair{
			Login:    r.PostForm.Get("login"),
			Password: r.PostForm.Get("password"),
		},
		Metainfo: model.Metainfo{
			Date: time.Now(),
			Text: r.PostForm.Get("meta"),
		},
	}
	err := h.store.AddPair(r.Context(), userID, pairs)
	if err != nil {
		h.ErrorWithLog(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetPairs(w http.ResponseWriter, r *http.Request) {
	rCtx := r.Context()
	userID := rCtx.Value(model.UserIDCtxKey{}).(int)
	limit := rCtx.Value(model.LimitCtxKey{}).(int)
	offset := rCtx.Value(model.OffsetCtxKey{}).(int)

	pairs, err := h.store.GetPairs(rCtx, userID, limit, offset)
	if err != nil {
		h.ErrorWithLog(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(pairs)
	if err != nil {
		h.ErrorWithLog(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(resp)
	if err != nil {
		h.ErrorWithLog(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
