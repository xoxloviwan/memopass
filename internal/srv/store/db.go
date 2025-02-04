package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"iwakho/gopherkeep/internal/model"
	"iwakho/gopherkeep/internal/srv/errs"
	"mime/multipart"
	"time"

	"github.com/ncruces/go-sqlite3"
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

func (db *Storage) AddPair(ctx context.Context, userID int, pair model.PairInfo) error {
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

func (db *Storage) AddFile(ctx context.Context, userID int, file io.Reader, fh *multipart.FileHeader, isBinary bool) error {
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
		old := db.SetInterrupt(ctx)
		defer db.SetInterrupt(old)

		const (
			prefix    = "@"
			userTag   = "user_id"
			dateTag   = "date"
			nameTag   = "name"
			sizeTag   = "size"
			binaryTag = "binary"
			fileTag   = "file"
		)
		query := fmt.Sprintf(`INSERT INTO files (%s, %s, %s, %s, %s) VALUES (%s, %s, %s, zeroblob(%s), %s)`,
			userTag, dateTag, nameTag, fileTag, binaryTag,
			prefix+userTag, prefix+dateTag, prefix+nameTag, prefix+sizeTag, prefix+binaryTag)

		stmt, _, err := db.Prepare(query)
		if err != nil {
			return err
		}

		err = stmt.BindInt(stmt.BindIndex(prefix+userTag), userID)
		if err != nil {
			return err
		}
		err = stmt.BindTime(stmt.BindIndex(prefix+dateTag), time.Now(), sqlite3.TimeFormatDefault)
		if err != nil {
			return err
		}
		err = stmt.BindText(stmt.BindIndex(prefix+nameTag), fh.Filename)
		if err != nil {
			return err
		}
		err = stmt.BindBool(stmt.BindIndex(prefix+binaryTag), isBinary)
		if err != nil {
			return err
		}
		err = stmt.BindInt64(stmt.BindIndex(prefix+sizeTag), fh.Size)
		if err != nil {
			return err
		}

		err = stmt.Exec()
		if err != nil {
			return err
		}
		err = stmt.ClearBindings()
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

func (db *Storage) GetFiles(ctx context.Context, userID int, limit int, offset int, isBinary bool) ([]model.FileInfo, error) {
	files := []model.FileInfo{}
	rows, err := db.QueryContext(ctx, `SELECT
				name,
				file,
				date,
				meta,
				id
			FROM files
			WHERE user_id = @user_id and binary = @binary ORDER BY date DESC LIMIT @limit OFFSET @offset`,
		sql.Named("user_id", userID),
		sql.Named("limit", limit),
		sql.Named("offset", offset),
		sql.Named("binary", isBinary),
	)
	if err != nil {
		return files, err
	}
	defer rows.Close()
	for rows.Next() {
		file := model.FileInfo{}
		var meta sql.NullString
		err = rows.Scan(&file.Name, &file.Blob, &file.Date, &meta, &file.ID)
		if err != nil {
			return files, err
		}
		if meta.Valid {
			file.Text = meta.String
		}
		files = append(files, file)
	}
	return files, nil
}

func (db *Storage) AddCard(ctx context.Context, userID int, card model.CardInfo) error {
	_, err := db.ExecContext(ctx, "INSERT INTO cards (user_id, ccn, exp, cvv, date, meta) VALUES ($1, $2, $3, $4, $5, $6)", userID, card.Number, card.Exp, card.VerifVal, card.Date, card.Text)
	return err
}

func (db *Storage) GetCards(ctx context.Context, userID int, limit int, offset int) ([]model.CardInfo, error) {
	cards := []model.CardInfo{}
	rows, err := db.QueryContext(ctx, `SELECT ccn, exp, cvv, date, meta
			FROM cards
			WHERE user_id = @user_id ORDER BY date DESC LIMIT @limit OFFSET @offset`,
		sql.Named("user_id", userID),
		sql.Named("limit", limit),
		sql.Named("offset", offset),
	)
	if err != nil {
		return cards, err
	}
	defer rows.Close()

	for rows.Next() {
		card := model.CardInfo{}
		err = rows.Scan(&card.Number, &card.Exp, &card.VerifVal, &card.Date, &card.Text)
		if err != nil {
			return cards, err
		}
		cards = append(cards, card)
	}

	return cards, nil
}

func (db *Storage) GetFileById(ctx context.Context, userID int, id int, isBinary bool) (data []byte, name string, err error) {
	err = db.QueryRowContext(ctx, "SELECT file, name FROM files WHERE user_id = @user_id and id = @id and binary = @binary", sql.Named("user_id", userID), sql.Named("id", id), sql.Named("binary", isBinary)).Scan(&data, &name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", errs.ErrNotFound
		}
		return nil, "", err
	}
	return data, name, nil
}
