package controls

import (
	"bytes"
	"encoding/json"
	"fmt"
	iHttp "iwakho/gopherkeep/internal/cli/http"
	"iwakho/gopherkeep/internal/model"
	"net/http"
)

var (
	token string
)

func TryLogin(login string, password string) error {
	loginUrl := iHttp.BaseURL + "/user/login"
	r, err := http.NewRequest("GET", loginUrl, nil)
	if err != nil {
		return fmt.Errorf("Bad request: url=%s e=%s", loginUrl, err)
	}
	r.SetBasicAuth(login, password)
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
	return nil
}

func SignUp(login string, password string) error {
	signUpUrl := iHttp.BaseURL + "/user/signup"
	creds := model.Creds{
		User: login,
		Pwd:  password,
	}
	body, err := json.Marshal(creds)
	if err != nil {
		return err
	}
	r, err := http.NewRequest("POST", signUpUrl, bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")
	if err != nil {
		return fmt.Errorf("Bad request: url=%s e=%s", signUpUrl, err)
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
	return nil
}
