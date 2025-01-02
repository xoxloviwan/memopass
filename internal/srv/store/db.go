package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"iwakho/gopherkeep/internal/model"
	"mime/multipart"
	"time"

	drv "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
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

func (db *Storage) AddNewPair(ctx context.Context, userID int, pair model.PairInfo) error {
	_, err := db.ExecContext(ctx, "INSERT INTO pairs (user_id, date, login, password, meta) VALUES ($1, $2, $3, $4, $5)", userID, pair.Date, pair.Login, pair.Password, pair.Text)
	return err
}

func (db *Storage) GetPairs(ctx context.Context, userID int, limit int, offset int) ([]model.PairInfo, error) {
	pairs := []model.PairInfo{}
	rows, err := db.QueryContext(ctx, `SELECT
				login,
				password,
				date,
				meta
			FROM pairs
			WHERE user_id = @user_id ORDER BY date DESC LIMIT @limit OFFSET @offset`,
		sql.Named("user_id", userID),
		sql.Named("limit", limit),
		sql.Named("offset", offset),
	)
	if err != nil {
		return pairs, err
	}
	defer rows.Close()

	for rows.Next() {
		pair := model.PairInfo{}
		err = rows.Scan(&pair.Login, &pair.Password, &pair.Date, &pair.Text)
		if err != nil {
			return pairs, err
		}
		pairs = append(pairs, pair)
	}

	return pairs, nil
}

func (db *Storage) AddFile(ctx context.Context, userID int, file io.Reader, fh *multipart.FileHeader) error {
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	return conn.Raw(func(driverConn any) error {
		conn, ok := driverConn.(drv.Conn)
		if !ok {
			return errors.New("not driver.Conn")
		}
		db := conn.Raw()
		defer db.Close()
		db.SetInterrupt(ctx)

		err = db.Exec(fmt.Sprintf(`INSERT INTO files (user_id, date, name, file) VALUES ('%d', '%s', '%s', zeroblob(%d))`, userID, time.Now(), fh.Filename, fh.Size))
		if err != nil {
			return err
		}
		blob, err := db.OpenBlob("main", "files", "file", db.LastInsertRowID(), true)
		if err != nil {
			return err
		}
		defer blob.Close()
		n, err := io.Copy(blob, file)
		if err != nil {
			return err
		}
		fmt.Printf("written %d\n", n)

		return nil
	})
}
