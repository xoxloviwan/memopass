package store

import (
	"context"
	"errors"
	"iwakho/gopherkeep/internal/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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
