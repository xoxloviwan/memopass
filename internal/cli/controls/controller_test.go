package controls

import (
	iCrypto "iwakho/gopherkeep/internal/cli/crypto"
	iHttp "iwakho/gopherkeep/internal/cli/http"
	"iwakho/gopherkeep/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func setup(t *testing.T, baseURL string, opts ...func(*Controller)) *Controller {
	cli, err := iHttp.New("", baseURL)
	if err != nil {
		t.Fatal(err)
	}
	ctrl := New(cli)
	for _, opt := range opts {
		opt(ctrl)
	}

	return ctrl
}
func WithEncryption() func(*Controller) {
	return func(c *Controller) {
		c.crypto = iCrypto.NewCryptoManager(model.Pair{Login: "alice", Password: "12345"})
	}
}

func TestLoginCli(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		if diff := cmp.Diff(req.URL.String(), "/api/v1/user/login"); diff != "" {
			t.Error(diff)
		}
		rw.Header().Set("Authorization", "token")
		rw.WriteHeader(http.StatusOK)
	}))
	// Close the server when test finishes
	defer server.Close()

	// Use Client & URL from our local test server

	c := setup(t, server.URL)
	err := c.Login(model.Pair{Login: "alice", Password: "12345"})

	if err != nil {
		t.Error(err)
		return
	}
}

func TestAddCardCli(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		if diff := cmp.Diff(req.URL.String(), "/api/v1/item/add/card"); diff != "" {
			t.Error(diff)
		}
		rw.WriteHeader(http.StatusOK)
	}))
	// Close the server when test finishes
	defer server.Close()

	// Use Client & URL from our local test server

	c := setup(t, server.URL, WithEncryption())
	err := c.AddCard(model.Card{Number: "1234567890123456", Exp: "12/25", VerifVal: "123"})
	if err != nil {
		t.Error(err)
		return
	}
}
