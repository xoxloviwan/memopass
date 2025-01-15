package http

import (
	"bytes"
	"iwakho/gopherkeep/internal/srv/http/handlers"
	mock "iwakho/gopherkeep/internal/srv/http/handlers/mockstore"
	"iwakho/gopherkeep/internal/srv/jwt"
	iLog "iwakho/gopherkeep/internal/srv/log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

type want struct {
	code        int
	contentType string
	err         error
}

type testcase struct {
	name   string
	url    string
	method string
	want   want
	pars   map[string]string
}

func setup(t *testing.T) (http.Handler, *mock.MockStore) {
	logger := iLog.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := mock.NewMockStore(ctrl)
	hdl := handlers.NewHandler(db, logger)
	router := NewRouter(http.NewServeMux(), logger)
	return router.SetupRoutes(hdl), db
}

func multipartBody(t *testing.T, pars map[string]string) (*bytes.Buffer, string) {
	var err error
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	for k, v := range pars {
		err = w.WriteField(k, v)
		if err != nil {
			t.Error(err)
			return nil, ""
		}
	}
	err = w.Close()
	if err != nil {
		t.Error(err)
		return nil, ""
	}
	return body, w.FormDataContentType()
}

func Test_AddCard(t *testing.T) {
	tests := []testcase{
		{
			name:   "success",
			url:    "/api/v1/item/add?type=3",
			method: http.MethodPost,
			pars: map[string]string{
				"cnn": "1234 5678 9012 3456",
				"exp": "12/25",
				"cvv": "123",
			},
			want: want{
				code: http.StatusOK,
			},
		},
	}
	router, db := setup(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := 1
			tkn, err := jwt.BuildJWT("user", userID)
			if err != nil {
				t.Error(err)
			}
			body, contentTypeHeader := multipartBody(t, tt.pars)
			req := httptest.NewRequest(tt.method, tt.url, body)
			w := httptest.NewRecorder()

			db.EXPECT().AddCard(gomock.Any(), userID, gomock.Any()).Return(tt.want.err)
			req.Header.Set("Authorization", jwt.Bearer+tkn)
			req.Header.Set("Content-Type", contentTypeHeader)
			router.ServeHTTP(w, req)

			if w.Code != tt.want.code {
				t.Errorf("expected %v; got %v", tt.want.code, w.Code)
			}
		})
	}
}

func Test_AddPair(t *testing.T) {
	tests := []testcase{
		{
			name:   "success",
			url:    "/api/v1/item/add?type=0",
			method: http.MethodPost,
			pars: map[string]string{
				"login":    "vasya",
				"password": "123456",
			},
			want: want{
				code: http.StatusOK,
			},
		},
	}
	router, db := setup(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := 1
			tkn, err := jwt.BuildJWT("user", userID)
			if err != nil {
				t.Error(err)
			}
			body, contentTypeHeader := multipartBody(t, tt.pars)
			req := httptest.NewRequest(tt.method, tt.url, body)
			w := httptest.NewRecorder()

			db.EXPECT().AddNewPair(gomock.Any(), userID, gomock.Any()).Return(tt.want.err)
			req.Header.Set("Authorization", jwt.Bearer+tkn)
			req.Header.Set("Content-Type", contentTypeHeader)
			router.ServeHTTP(w, req)

			if w.Code != tt.want.code {
				t.Errorf("expected %v; got %v", tt.want.code, w.Code)
			}
		})
	}
}
