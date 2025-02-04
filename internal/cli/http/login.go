package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"iwakho/gopherkeep/internal/model"
	"net/http"
)

func (cli *Client) Login(p model.Pair) error {
	r, err := http.NewRequest("GET", cli.Api.login, nil)
	if err != nil {
		return fmt.Errorf("Bad request: url=%s e=%s", cli.Api.login, err)
	}
	r.SetBasicAuth(p.Login, p.Password)
	resp, err := cli.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Bad status code: %d url=%s", resp.StatusCode, cli.Api.login)
	}
	token := resp.Header.Get("Authorization")
	if token == "" {
		return fmt.Errorf("bad token")

	}
	cli.token = token
	return nil
}

func (cl *Client) SignUp(p model.Pair) error {
	creds := model.Creds{
		User: p.Login,
		Pwd:  p.Password,
	}
	body, err := json.Marshal(creds)
	if err != nil {
		return err
	}
	r, err := http.NewRequest("POST", cl.Api.signUp, bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")
	if err != nil {
		return fmt.Errorf("Bad request: url=%s e=%s", cl.Api.signUp, err)
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
	return nil
}
