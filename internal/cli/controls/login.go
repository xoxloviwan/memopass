package controls

import (
	"fmt"
	iHttp "iwakho/gopherkeep/internal/cli/http"
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
