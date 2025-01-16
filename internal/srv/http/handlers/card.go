package handlers

import (
	"encoding/json"
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

func (h *Handler) GetCards(w http.ResponseWriter, r *http.Request) {
	rCtx := r.Context()
	userID := rCtx.Value(model.UserIDCtxKey{}).(int)
	limit := rCtx.Value(model.LimitCtxKey{}).(int)
	offset := rCtx.Value(model.OffsetCtxKey{}).(int)

	cards, err := h.store.GetCards(rCtx, userID, limit, offset)
	if err != nil {
		h.ErrorWithLog(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(cards)
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
