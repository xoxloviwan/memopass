package handlers

import (
	"encoding/json"
	"iwakho/gopherkeep/internal/model"
	"net/http"
	"strconv"
	"time"
)

const maxMemory = 5 << 20    // 5 MB
const maxFileSize = 10 << 20 // 10 MB

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
			Metainfo: model.Metainfo{
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
		fhs := r.MultipartForm.File["file"]
		if len(fhs) == 0 {
			http.Error(w, "no file provided", http.StatusBadRequest)
			return
		}
		f0 := fhs[0]
		if f0.Size > maxFileSize {
			http.Error(w, "file too large", http.StatusBadRequest)
			return
		}
		file, err := f0.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		err = h.store.AddFile(r.Context(), userID, file, f0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	case model.ItemTypeCard:
	// TODO: add card item
	default:
		http.Error(w, "unknown item type", http.StatusBadRequest)
	}
}

func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.UserIDCtxKey{}).(int)
	queries := r.URL.Query()
	itemTypeStr := queries.Get("type")
	limitStr := queries.Get("limit")
	offsetStr := queries.Get("offset")

	itemType, err := strconv.Atoi(itemTypeStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var limit, offset int
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	switch itemType {
	case model.ItemTypeLoginPass:
		pairs, err := h.store.GetPairs(r.Context(), userID, limit, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp, err := json.Marshal(pairs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
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
