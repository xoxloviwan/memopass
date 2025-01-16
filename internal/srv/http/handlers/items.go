package handlers

import (
	"iwakho/gopherkeep/internal/model"
	"net/http"
	"strconv"
)

func (h *Handler) GetItems(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.UserIDCtxKey{}).(int)
	queries := r.URL.Query()
	itemTypeStr := queries.Get("type")
	limitStr := queries.Get("limit")
	offsetStr := queries.Get("offset")

	itemType, err := strconv.Atoi(itemTypeStr)
	if err != nil {
		h.ErrorWithLog(w, err.Error(), http.StatusBadRequest)
		return
	}
	var limit, offset int
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			h.ErrorWithLog(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			h.ErrorWithLog(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	switch itemType {
	case model.ItemTypeLoginPass:
		h.getPairs(w, r, userID, limit, offset)
	case model.ItemTypeText:
	// TODO: add text item
	case model.ItemTypeBinary:
		h.getFiles(w, r, userID, limit, offset)
	case model.ItemTypeCard:
		h.getCards(w, r, userID, limit, offset)
	default:
		h.ErrorWithLog(w, "unknown item type", http.StatusBadRequest)
	}
}
