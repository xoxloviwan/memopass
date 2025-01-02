package handlers

import (
	"context"
	"io"
	"iwakho/gopherkeep/internal/model"
	"mime/multipart"
)

type Store interface {
	NewUser(ctx context.Context, user model.User) (int, error)
	GetUser(ctx context.Context, username string) (*model.User, error)
	AddNewPair(ctx context.Context, userID int, pair model.PairInfo) error
	GetPairs(ctx context.Context, userID int, limit int, offset int) ([]model.PairInfo, error)
	AddFile(ctx context.Context, userID int, file io.Reader, fh *multipart.FileHeader) error
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
