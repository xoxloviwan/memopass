package store

import (
	"context"
	"errors"
	"iwakho/gopherkeep/internal/model"
	"iwakho/gopherkeep/internal/srv/errs"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
)

func NewMockStorage() (*Storage, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	return &Storage{db}, mock, nil
}

func TestNewUser(t *testing.T) {

	cases := []struct {
		name       string
		wantErr    error
		wantUserID int
	}{
		{
			name:       "success",
			wantErr:    nil,
			wantUserID: 123,
		},
		{
			name:       "fail",
			wantErr:    errors.New("some error"),
			wantUserID: 0,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := NewMockStorage()
			if err != nil {
				t.Error(err)
			}
			defer db.Close()

			u := model.User{
				Name: "vasya",
				Hash: []byte("123456"),
			}

			if tt.wantErr != nil {
				mock.ExpectQuery("INSERT INTO users").WithArgs(u.Name, u.Hash).WillReturnError(tt.wantErr)
			} else {
				row := sqlmock.NewRows([]string{"id"}).AddRow(tt.wantUserID)
				mock.ExpectQuery("INSERT INTO users").WithArgs(u.Name, u.Hash).WillReturnRows(row)
			}
			gotUserID, err := db.NewUser(context.Background(), u)
			if err != tt.wantErr {
				t.Error(err)
				return
			}
			if err != nil {
				return
			}
			if tt.wantUserID != gotUserID {
				t.Errorf("want %d, got %d", tt.wantUserID, gotUserID)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	cases := []struct {
		name       string
		wantErr    error
		wantUserID int
	}{
		{
			name:       "success",
			wantErr:    nil,
			wantUserID: 123,
		},
		{
			name:       "fail",
			wantErr:    errors.New("some error"),
			wantUserID: 0,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := NewMockStorage()
			if err != nil {
				t.Error(err)
			}
			defer db.Close()

			if tt.wantErr != nil {
				mock.ExpectQuery("SELECT id, username, password FROM user").WithArgs("vasya").WillReturnError(tt.wantErr)
			} else {
				row := sqlmock.NewRows([]string{"id", "username", "password"}).AddRow(tt.wantUserID, "vasya", "123456")
				mock.ExpectQuery("SELECT id, username, password FROM user").WithArgs("vasya").WillReturnRows(row)
			}
			gotUserID, err := db.GetUser(context.Background(), "vasya")
			if err != tt.wantErr {
				t.Error(err)
				return
			}
			if err != nil {
				return
			}
			if tt.wantUserID != gotUserID.ID {
				t.Errorf("want %d, got %d", tt.wantUserID, gotUserID.ID)
			}
		})
	}
}

func TestAddPair(t *testing.T) {
	cases := []struct {
		name    string
		wantErr error
	}{
		{
			name:    "success",
			wantErr: nil,
		},
		{
			name:    "fail",
			wantErr: errors.New("some error"),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := NewMockStorage()
			if err != nil {
				t.Error(err)
			}
			defer db.Close()

			userID := 1
			pair := model.PairInfo{
				Pair: model.Pair{
					Login:    "vasya",
					Password: "123456",
				},
				Metainfo: model.Metainfo{
					Date: time.Now(),
					Text: "some meta",
				},
			}

			mock.ExpectExec("INSERT INTO pairs").WithArgs(userID, pair.Date, pair.Login, pair.Password, pair.Text).WillReturnError(tt.wantErr).WillReturnResult(sqlmock.NewResult(1, 1))
			err = db.AddPair(context.Background(), userID, pair)
			if err != tt.wantErr {
				t.Error(err)
				return
			}
		})
	}
}

func TestAddCard(t *testing.T) {
	cases := []struct {
		name    string
		wantErr error
	}{
		{
			name:    "success",
			wantErr: nil,
		},
		{
			name:    "fail",
			wantErr: errors.New("some error"),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := NewMockStorage()
			if err != nil {
				t.Error(err)
			}
			defer db.Close()

			userID := 1
			card := model.CardInfo{
				Card: model.Card{
					Number:   "1234567890123456",
					Exp:      "12/22",
					VerifVal: "123",
				},
				Metainfo: model.Metainfo{
					Date: time.Now(),
					Text: "some meta",
				},
			}

			mock.ExpectExec("INSERT INTO cards").WithArgs(userID, card.Number, card.Exp, card.VerifVal, card.Date, card.Text).WillReturnError(tt.wantErr).WillReturnResult(sqlmock.NewResult(1, 1))
			err = db.AddCard(context.Background(), userID, card)
			if err != tt.wantErr {
				t.Error(err)
				return
			}
		})
	}
}

func TestGetCards(t *testing.T) {
	cases := []struct {
		name    string
		wantErr error
	}{
		{
			name:    "success",
			wantErr: nil,
		},
		{
			name:    "fail",
			wantErr: errors.New("some error"),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := NewMockStorage()
			if err != nil {
				t.Error(err)
			}
			defer db.Close()

			userID := 1
			cards := []model.CardInfo{
				{
					Card: model.Card{
						Number:   "1234567890123456",
						Exp:      "12/22",
						VerifVal: "123",
					},
					Metainfo: model.Metainfo{
						Date: time.Now(),
						Text: "some meta",
					},
				},
			}
			rows := sqlmock.NewRows([]string{"card_number", "exp", "verif_val", "date", "meta"}).
				AddRow(cards[0].Number, cards[0].Exp, cards[0].VerifVal, cards[0].Date, cards[0].Text)
			mock.ExpectQuery("SELECT").WithArgs(userID, 10, 0).WillReturnRows(rows).WillReturnError(tt.wantErr)
			got, err := db.GetCards(context.Background(), userID, 10, 0)
			if err != tt.wantErr {
				t.Error(err)
				return
			}
			if tt.wantErr != nil {
				return
			}

			if err != nil {
				t.Error(err)
				return
			}
			if len(got) != len(cards) {
				t.Errorf("want %d, got %d", len(cards), len(got))
			}
		})
	}
}

func TestGetPairs(t *testing.T) {
	cases := []struct {
		name    string
		wantErr error
	}{
		{
			name:    "success",
			wantErr: nil,
		},
		{
			name:    "fail",
			wantErr: errors.New("some error"),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := NewMockStorage()
			if err != nil {
				t.Error(err)
			}
			defer db.Close()

			userID := 1
			pairs := []model.PairInfo{
				{
					Pair: model.Pair{
						Login:    "vasya",
						Password: "123456",
					},
					Metainfo: model.Metainfo{
						Date: time.Now(),
						Text: "some meta",
					},
				},
				{
					Pair: model.Pair{
						Login:    "petya",
						Password: "5678",
					},
					Metainfo: model.Metainfo{
						Date: time.Now(),
						Text: "some meta",
					},
				},
			}
			rows := sqlmock.NewRows([]string{"login", "password", "date", "meta"}).
				AddRow(pairs[0].Login, pairs[0].Password, pairs[0].Date, pairs[0].Text).
				AddRow(pairs[1].Login, pairs[1].Password, pairs[1].Date, pairs[1].Text)
			mock.ExpectQuery("SELECT").WithArgs(userID, 10, 0).WillReturnRows(rows).WillReturnError(tt.wantErr)
			got, err := db.GetPairs(context.Background(), userID, 10, 0)
			if err != tt.wantErr {
				t.Error(err)
				return
			}
			if tt.wantErr != nil {
				return
			}

			if err != nil {
				t.Error(err)
				return
			}
			if len(got) != len(pairs) {
				t.Errorf("want %d, got %d", len(pairs), len(got))
			}
		})
	}
}

func TestGetFiles(t *testing.T) {
	cases := []struct {
		name    string
		wantErr error
	}{
		{
			name:    "success",
			wantErr: nil,
		},
		{
			name:    "fail",
			wantErr: errors.New("some error"),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := NewMockStorage()
			if err != nil {
				t.Error(err)
			}
			defer db.Close()

			userID := 1
			files := []model.FileInfo{
				{
					File: model.File{
						Name: "file1",
						Blob: []byte("blob1"),
						ID:   1,
					},
					Metainfo: model.Metainfo{
						Date: time.Now(),
						Text: "some meta",
					},
				},
				{
					File: model.File{
						Name: "file2",
						Blob: []byte("blob2"),
						ID:   2,
					},
					Metainfo: model.Metainfo{
						Date: time.Now(),
						Text: "some meta",
					},
				},
			}
			getRows := func() *sqlmock.Rows {
				return sqlmock.NewRows([]string{"name", "path", "date", "meta", "id"}).
					AddRow(files[0].Name, files[0].Blob, files[0].Date, files[0].Text, files[0].ID).
					AddRow(files[1].Name, files[1].Blob, files[1].Date, files[1].Text, files[1].ID)
			}
			rows := getRows()
			t.Log(tt.name, 1)
			mock.ExpectQuery("SELECT").WithArgs(userID, 10, 0, false).WillReturnRows(rows).WillReturnError(tt.wantErr)
			got, err := db.GetFiles(context.Background(), userID, 10, 0, false)
			if err != tt.wantErr {
				t.Error(err)
				return
			}
			if tt.wantErr != nil {
				return
			}

			if err != nil {
				t.Error(err)
				return
			}
			if len(got) != len(files) {
				t.Errorf("want %d, got %d", len(files), len(got))
			}
			t.Log(tt.name, 2)
			rows = getRows()
			mock.ExpectQuery("SELECT").WithArgs(userID, 10, 0, true).WillReturnRows(rows).WillReturnError(tt.wantErr)
			got, err = db.GetFiles(context.Background(), userID, 10, 0, true)
			if err != tt.wantErr {
				t.Error(err)
				return
			}
			if tt.wantErr != nil {
				return
			}
			if err != nil {
				t.Error(err)
				return
			}
			if len(got) != len(files) {
				t.Errorf("want %d, got %d", len(files), len(got))
			}
		})
	}
}

func TestGetFilesByID(t *testing.T) {
	cases := []struct {
		name    string
		wantErr error
		file    model.File
	}{
		{
			name:    "success",
			wantErr: nil,
			file: model.File{
				Name: "file1",
				Blob: []byte("blob1"),
				ID:   1,
			},
		},
		{
			name:    "fail",
			wantErr: errors.New("some error"),
			file:    model.File{},
		},
		{
			name:    "no file",
			wantErr: errs.ErrNotFound,
			file:    model.File{},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := NewMockStorage()
			if err != nil {
				t.Error(err)
			}
			defer db.Close()

			userID := 1
			rows := sqlmock.NewRows([]string{"name", "file"}).
				AddRow(tt.file.Blob, tt.file.Name)

			mock.ExpectQuery("SELECT").WithArgs(userID, tt.file.ID, false).WillReturnRows(rows).WillReturnError(tt.wantErr)
			gotBytes, gotName, err := db.GetFileById(context.Background(), userID, tt.file.ID, false)
			if err != tt.wantErr {
				t.Error(err)
				return
			}
			if tt.wantErr != nil {
				return
			}
			if err != nil {
				t.Error(err)
				return
			}

			if diff := cmp.Diff(gotBytes, tt.file.Blob); diff != "" {
				t.Errorf("Body mismatch (-want +got):\n%s", diff)
				return
			}

			if diff := cmp.Diff(gotName, tt.file.Name); diff != "" {
				t.Errorf("Name mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}
