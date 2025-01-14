package controls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"iwakho/gopherkeep/internal/cli/crypto"
	iHttp "iwakho/gopherkeep/internal/cli/http"
	"iwakho/gopherkeep/internal/model"
	"net/http"
)

var (
	token    string
	CrptMngr *crypto.CryptoManager
)

type LoginCtrl struct{}

func (LoginCtrl) Submit(p model.Pair) error {
	r, err := http.NewRequest("GET", iHttp.ApiLogin, nil)
	if err != nil {
		return fmt.Errorf("Bad request: url=%s e=%s", iHttp.ApiLogin, err)
	}
	r.SetBasicAuth(p.Login, p.Password)
	resp, err := iHttp.Client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Bad status code: %d", resp.StatusCode)
	}
	token = resp.Header.Get("Authorization")
	if token == "" {
		return fmt.Errorf("bad token")

	}
	CrptMngr = crypto.NewCryptoManager(p)
	return nil
}

type SignUpCtrl struct{}

func (SignUpCtrl) Submit(p model.Pair) error {
	creds := model.Creds{
		User: p.Login,
		Pwd:  p.Password,
	}
	body, err := json.Marshal(creds)
	if err != nil {
		return err
	}
	r, err := http.NewRequest("POST", iHttp.ApiSignUp, bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")
	if err != nil {
		return fmt.Errorf("Bad request: url=%s e=%s", iHttp.ApiSignUp, err)
	}
	resp, err := iHttp.Client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Bad status code: %d", resp.StatusCode)
	}
	token = resp.Header.Get("Authorization")
	if token == "" {
		return fmt.Errorf("Bad token")
	}
	CrptMngr = crypto.NewCryptoManager(p)
	return nil
}
