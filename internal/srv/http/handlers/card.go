package handlers

import (
	"iwakho/gopherkeep/internal/model"
	"net/http"
	"time"
)

func (h *Handler) AddCard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.UserIDCtxKey{}).(int)
	card := model.CardInfo{
		Card: model.Card{
			Number:   r.PostForm.Get("ccn"),
			Exp:      r.PostForm.Get("exp"),
			VerifVal: r.PostForm.Get("cvv"),
		},
		Metainfo: model.Metainfo{
			Date: time.Now(),
			Text: r.PostForm.Get("meta"),
		},
	}
	err := h.store.AddCard(r.Context(), userID, card)
	if err != nil {
		h.ErrorWithLog(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
