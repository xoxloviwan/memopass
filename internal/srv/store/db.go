package store

import (
	"context"
	"database/sql"
	"iwakho/gopherkeep/internal/srv/model"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	*sql.DB
}

func NewStorage(dsn string) (*Storage, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	return &Storage{db}, nil
}

func (db *Storage) NewUser(ctx context.Context, u model.User) (int, error) {
	row := db.QueryRowContext(ctx, "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", u.Name, u.Hash)
	err := row.Scan(&u.ID)
	if err != nil {
		return 0, err
	}
	return u.ID, nil
}

func (db *Storage) GetUser(ctx context.Context, login string) (*model.User, error) {
	u := model.User{}
	row := db.QueryRowContext(ctx, "SELECT id, username, password FROM users WHERE username = $1", login)
	err := row.Scan(&u.ID, &u.Name, &u.Hash)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
