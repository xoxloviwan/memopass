package handlers

import (
	"iwakho/gopherkeep/internal/model"
	"net/http"
	"strconv"
)

const maxMemory = 5 << 20    // 5 MB
const maxFileSize = 10 << 20 // 10 MB

func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		h.ErrorWithLog(w, err.Error(), http.StatusBadRequest)
		return
	}
	itemTypeStr := r.URL.Query().Get("type")
	itemType, err := strconv.Atoi(itemTypeStr)
	if err != nil {
		h.ErrorWithLog(w, err.Error(), http.StatusBadRequest)
		return
	}
	switch itemType {
	case model.ItemTypeLoginPass:
		h.AddPair(w, r)
	case model.ItemTypeText:
	// TODO: add text item
	case model.ItemTypeBinary:
		h.AddFile(w, r)
	case model.ItemTypeCard:
		h.AddCard(w, r)
	default:
		h.ErrorWithLog(w, "unknown item type", http.StatusBadRequest)
	}
}

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
		// TODO: add card item
	default:
		h.ErrorWithLog(w, "unknown item type", http.StatusBadRequest)
	}
}
