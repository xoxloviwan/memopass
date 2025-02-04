package handlers

import (
	"context"
	"io"
	"iwakho/gopherkeep/internal/model"
	"mime/multipart"
	"net/http"
)

//go:generate mockgen -destination ./mockstore/mock_store.go -package mockstore . Store

type Store interface {
	NewUser(ctx context.Context, user model.User) (int, error)
	GetUser(ctx context.Context, username string) (*model.User, error)
	AddPair(ctx context.Context, userID int, pair model.PairInfo) error
	GetPairs(ctx context.Context, userID int, limit int, offset int) ([]model.PairInfo, error)
	AddFile(ctx context.Context, userID int, file io.Reader, fh *multipart.FileHeader, isBinary bool) error
	AddCard(ctx context.Context, userID int, card model.CardInfo) error
	GetFiles(ctx context.Context, userID int, limit int, offset int, isBinary bool) ([]model.FileInfo, error)
	GetFileById(ctx context.Context, userID int, id int, isBinary bool) (data []byte, name string, err error)
	GetCards(ctx context.Context, userID int, limit int, offset int) ([]model.CardInfo, error)
}

type logger interface {
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}

type Handler struct {
	store  Store
	logger logger
}

func NewHandler(store Store, logger logger) *Handler {
	return &Handler{
		store,
		logger,
	}
}

func (h *Handler) ErrorWithLog(w http.ResponseWriter, err string, code int) {
	h.logger.Error(err)
	http.Error(w, err, code)
}
