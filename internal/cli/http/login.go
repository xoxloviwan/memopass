package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"iwakho/gopherkeep/internal/cli/crypto"
	"iwakho/gopherkeep/internal/model"
	"net/http"
)

var (
	CrptMngr *crypto.CryptoManager
)

func (cli *Client) Login(p model.Pair) error {
	r, err := http.NewRequest("GET", cli.api.login, nil)
	if err != nil {
		return fmt.Errorf("Bad request: url=%s e=%s", cli.api.login, err)
	}
	r.SetBasicAuth(p.Login, p.Password)
	resp, err := cli.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Bad status code: %d url=%s", resp.StatusCode, cli.api.login)
	}
	token := resp.Header.Get("Authorization")
	if token == "" {
		return fmt.Errorf("bad token")

	}
	cli.token = token
	CrptMngr = crypto.NewCryptoManager(p)
	return nil
}

type SignUpCtrl struct{}

func (cl *Client) SignUp(p model.Pair) error {
	creds := model.Creds{
		User: p.Login,
		Pwd:  p.Password,
	}
	body, err := json.Marshal(creds)
	if err != nil {
		return err
	}
	r, err := http.NewRequest("POST", cl.api.signUp, bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")
	if err != nil {
		return fmt.Errorf("Bad request: url=%s e=%s", cl.api.signUp, err)
	}
	resp, err := cl.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Bad status code: %d", resp.StatusCode)
	}
	token := resp.Header.Get("Authorization")
	if token == "" {
		return fmt.Errorf("bad token")
	}
	cl.token = token
	CrptMngr = crypto.NewCryptoManager(p)
	return nil
}
