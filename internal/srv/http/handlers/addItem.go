package handlers

import (
	"iwakho/gopherkeep/internal/model"
	"net/http"
	"strconv"
	"time"
)

const maxMemory = 32 << 20 // 32 MB

func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.UserIDCtxKey{}).(int)
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	itemTypeStr := r.URL.Query().Get("type")
	itemType, err := strconv.Atoi(itemTypeStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	switch itemType {
	case model.ItemTypeLoginPass:
		pairs := model.PairInfo{
			Pair: model.Pair{
				Login:    r.PostForm.Get("login"),
				Password: r.PostForm.Get("password"),
			},
			Meta: model.Metainfo{
				Date: time.Now(),
				Text: r.PostForm.Get("meta"),
			},
		}
		err := h.store.AddNewPair(r.Context(), userID, pairs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	case model.ItemTypeText:
	// TODO: add text item
	case model.ItemTypeBinary:
	// TODO: add binary item
	case model.ItemTypeCard:
	// TODO: add card item
	default:
		http.Error(w, "unknown item type", http.StatusBadRequest)
	}
}
