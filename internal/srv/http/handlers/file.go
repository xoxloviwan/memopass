package handlers

import (
	"encoding/json"
	"iwakho/gopherkeep/internal/model"
	"net/http"
)

const (
	maxFileSize     = 10 << 20 // 10 MB
	errNoFile       = "no file provided"
	errFileTooLarge = "file too large"
)

func (h *Handler) AddBinary(w http.ResponseWriter, r *http.Request) {
	h.addFile(w, r, true)
}

func (h *Handler) AddText(w http.ResponseWriter, r *http.Request) {
	h.addFile(w, r, false)
}

func (h *Handler) addFile(w http.ResponseWriter, r *http.Request, isBinary bool) {
	userID := r.Context().Value(model.UserIDCtxKey{}).(int)
	fieldName := "file"
	if !isBinary {
		fieldName = "text"
	}
	fhs := r.MultipartForm.File[fieldName]
	if len(fhs) == 0 {
		h.ErrorWithLog(w, errNoFile, http.StatusBadRequest)
		return
	}
	f0 := fhs[0]
	if f0.Size > maxFileSize {
		h.ErrorWithLog(w, errFileTooLarge, http.StatusBadRequest)
		return
	}
	file, err := f0.Open()
	if err != nil {
		h.ErrorWithLog(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	err = h.store.AddFile(r.Context(), userID, file, f0, isBinary)
	if err != nil {
		h.ErrorWithLog(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetBinaries(w http.ResponseWriter, r *http.Request) {
	h.getFiles(w, r, true)
}

func (h *Handler) GetTexts(w http.ResponseWriter, r *http.Request) {
	h.getFiles(w, r, false)
}

func (h *Handler) getFiles(w http.ResponseWriter, r *http.Request, isBinary bool) {
	rCtx := r.Context()
	userID := rCtx.Value(model.UserIDCtxKey{}).(int)
	limit := rCtx.Value(model.LimitCtxKey{}).(int)
	offset := rCtx.Value(model.OffsetCtxKey{}).(int)

	files, err := h.store.GetFiles(rCtx, userID, limit, offset, isBinary)
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
