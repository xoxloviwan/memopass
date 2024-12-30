package controls

import (
	"bytes"
	"fmt"
	"io"
	iHttp "iwakho/gopherkeep/internal/cli/http"
	"iwakho/gopherkeep/internal/model"
	"mime/multipart"
	"net/http"
)

type AddPairCtrl struct {
	model.Pair
}

func (ctrl AddPairCtrl) Submit(login string, password string) error {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	err := w.WriteField("login", login)
	if err != nil {
		return err
	}
	err = w.WriteField("password", password)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	r, err := http.NewRequest("POST", iHttp.ApiAddItem, body)
	if err != nil {
		return err
	}
	r.Header.Set("Authorization", token)
	r.Header.Set("Content-Type", w.FormDataContentType())
	q := r.URL.Query()
	q.Add("type", "0")
	r.URL.RawQuery = q.Encode()
	resp, err := iHttp.Client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Bad status code: %d %s", resp.StatusCode, string(data))
	}
	return nil
}
