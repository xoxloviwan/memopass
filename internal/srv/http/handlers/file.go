package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"iwakho/gopherkeep/internal/model"
	"iwakho/gopherkeep/internal/srv/errs"
	"net/http"
	"strconv"
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

func (h *Handler) GetBinaryById(w http.ResponseWriter, r *http.Request) {
	h.getFileById(w, r, true)
}

func (h *Handler) GetTextById(w http.ResponseWriter, r *http.Request) {
	h.getFileById(w, r, false)
}

func (h *Handler) getFileById(w http.ResponseWriter, r *http.Request, isBinary bool) {
	rCtx := r.Context()
	userID := rCtx.Value(model.UserIDCtxKey{}).(int)
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.ErrorWithLog(w, "query param id is empty", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.ErrorWithLog(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, name, err := h.store.GetFileById(rCtx, userID, id, isBinary)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			h.ErrorWithLog(w, err.Error(), http.StatusNotFound)
			return
		}
		h.ErrorWithLog(w, err.Error(), http.StatusInternalServerError)
		return
	}
	contentType := "text/plain"
	if isBinary {
		contentType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", contentType)
	if name != "" {
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, name))
	}
	_, err = w.Write(file)
	if err != nil {
		h.ErrorWithLog(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
