package handlers

import (
	"encoding/json"
	"iwakho/gopherkeep/internal/model"
	"net/http"
)

func (h *Handler) AddFile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(model.UserIDCtxKey{}).(int)
	fhs := r.MultipartForm.File["file"]
	if len(fhs) == 0 {
		h.ErrorWithLog(w, "no file provided", http.StatusBadRequest)
		return
	}
	f0 := fhs[0]
	if f0.Size > maxFileSize {
		h.ErrorWithLog(w, "file too large", http.StatusBadRequest)
		return
	}
	file, err := f0.Open()
	if err != nil {
		h.ErrorWithLog(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	err = h.store.AddFile(r.Context(), userID, file, f0)
	if err != nil {
		h.ErrorWithLog(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getFiles(w http.ResponseWriter, r *http.Request, userID, limit, offset int) {
	files, err := h.store.GetFiles(r.Context(), userID, limit, offset)
	if err != nil {
		h.ErrorWithLog(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(files)
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
