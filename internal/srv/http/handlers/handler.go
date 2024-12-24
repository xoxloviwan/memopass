package handlers

import (
	"context"
	"iwakho/gopherkeep/internal/model"
)

type Store interface {
	NewUser(ctx context.Context, user model.User) (int, error)
	GetUser(ctx context.Context, username string) (*model.User, error)
	AddNewPair(ctx context.Context, userID int, pair model.Pairs) error
}

type logger interface {
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
